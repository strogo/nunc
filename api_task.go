package nunc

import (
	"fmt"
	"github.com/cznic/ql"
	"io/ioutil"
	"os"
	"time"
)

var (
	getTask          = ql.MustCompile("select * from Task where Context == $1 && EID == $2;")
	listTasks        = ql.MustCompile(fmt.Sprintf("select EID, Text, State from Task where Context == $1 && State > %d order by EID;", Resolved))
	addTask          = ql.MustCompile("insert into Task values($1, $2, $3, $4, $5, $6, $7);")
	delTask          = ql.MustCompile("delete from Task where Context == $1 && EID == $2;")
	getOldestTask    = ql.MustCompile("select * from Task where Context == $1 order by Modification asc limit 1;")
	getOldestTaskAny = ql.MustCompile("select * from Task order by Modification asc limit 1;")
	listAllTasks     = ql.MustCompile("select EID, Text, State from Task where Context == $1 order by EID;")
	doneTask         = ql.MustCompile("update Task set State = $1 where Context == $2 && EID == $3;")
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
		task.EID = id
		err = fmt.Errorf("task '%s' not found", TaskID(context, task))
		return
	}
	return
}

// TODO list -any (context)
func List(context Context, all bool) (tasks []Task, err error) {
	q := listTasks
	if all {
		q = listAllTasks
	}
	if err = query(nil, q, func(data []interface{}) (bool, error) {
		task := Task{
			EID:   data[0].(int64),
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

func Add(context Context, text, body string) (task Task, err error) {
	task = Task{
		Context:      context.ID,
		Text:         text,
		State:        Open,
		Creation:     time.Now(),
		Modification: time.Now(),
	}
	err = add(context, &task, body)
	return
}

func Delete(context Context, id int64) (task Task, err error) {
	task, err = Get(context, id, true)
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
	if err = execute(ctx, delTask, context.ID, id); err != nil {
		return
	}
	path := ResolvePath(TaskPath(context, task))
	if err = os.RemoveAll(path); err != nil {
		return
	}
	return
}

func add(context Context, task *Task, body string) (err error) {
	eid, err := GetNextEID(context.ShortName)
	if err != nil {
		return
	}
	task.EID = eid
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, addTask, ql.MustMarshal(*task)...); err != nil {
		return
	}
	if body != "" {
		bodyFile := ResolvePath(TaskPath(context, *task))
		err = ioutil.WriteFile(bodyFile, []byte(body), 0640)
		if err != nil {
			return
		}
	}
	return
}

func Fork(context Context, text, body string, parent int64) (task Task, err error) {
	task = Task{
		Context:      context.ID,
		Text:         text,
		State:        Open,
		Creation:     time.Now(),
		Modification: time.Now(),
		Parent:       parent,
	}
	err = add(context, &task, body)
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

func GetOldest(context Context) (task Task, err error) {
	if err = query(nil, getOldestTask, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&task, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}, context.ID); err != nil {
		return
	}
	if (task == Task{}) {
		err = fmt.Errorf("no task found")
		return
	}
	return
}

func GetOldestFromAllContexts() (task Task, context Context, err error) {
	if err = query(nil, getOldestTaskAny, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&task, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}); err != nil {
		return
	}
	if (task == Task{}) {
		err = fmt.Errorf("no task found")
		return
	}
	context, err = GetContextById(task.Context, true)
	if err != nil {
		return
	}
	return
}
