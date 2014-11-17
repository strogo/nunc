package nunc

import (
	"fmt"
	"github.com/cznic/ql"
	"io/ioutil"
	"time"
)

var (
	getTask      = ql.MustCompile("select * from Task where Context == $1 && id() == $2;")
	listTasks    = ql.MustCompile(fmt.Sprintf("select id(), Text, State from Task where Context == $1 && State > %d order by id();", Resolved))
	listAllTasks = ql.MustCompile("select id(), Text, State from Task where Context == $1 order by id();")
	doneTask     = ql.MustCompile("update Task set State = $1 where Context == $2 && id() == $3;")
	addTask      = ql.MustCompile("insert into Task values($1, $2, $3, $4, $5, $6);")
)

func Get(context Context, id int64, must bool) (task Task, err error) {
	if err = query(nil, getTask, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&task, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}, context.ID, id); err != nil {
		return
	}
	if (task == Task{}) && must {
		task.ID = id
		err = fmt.Errorf("task '%s' not found", TaskID(context, task))
	}
	return
}

func List(context Context, all bool) (tasks []Task, err error) {
	q := listTasks
	if all {
		q = listAllTasks
	}
	if err = query(nil, q, func(data []interface{}) (bool, error) {
		task := Task{
			ID:    data[0].(int64),
			Text:  data[1].(string),
			State: unmarshalStateFromInt(data[2].(int32)),
		}
		tasks = append(tasks, task)
		return true, nil
	}, context.ID); err != nil {
		return
	}
	return
}

func unmarshalStateFromInt(src int32) (state State) {
	switch src {
	case -1:
		state = Closed
	case 0:
		state = Resolved
	case 1:
		state = Open
	case 2:
		state = InProgress
	case 3:
		state = Feedback
	}
	return
}

func Add(context Context, text, body string) (err error) {
	task := Task{
		Context:      context.ID,
		Text:         text,
		State:        Open,
		Creation:     time.Now(),
		Modification: time.Now(),
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, addTask, ql.MustMarshal(task)...); err != nil {
		return
	}
	task.ID = ctx.LastInsertID
	if body != "" {
		bodyFile := ResolvePath(TaskPath(context, task))
		err = ioutil.WriteFile(bodyFile, []byte(body), 0640)
		if err != nil {
			return
		}
	}
	return
}

func Done(context Context, id int64, closed bool) (err error) {
	_, err = Get(context, id, true)
	if err != nil {
		return
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	state := Resolved
	if closed {
		state = Closed
	}
	if err = execute(ctx, doneTask, int32(state), context.ID, id); err != nil {
		return
	}
	return
}
