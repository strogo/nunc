package nunc

import (
	"fmt"
	"github.com/cznic/ql"
)

type Lock struct {
	Context int64 `ql:"index xContext"`
	Task    int64 `ql:"index xTask"`
}

var (
	getLock    = ql.MustCompile("select * from Lock where Context == $1;")
	setLockIn  = ql.MustCompile("insert into Lock values($1, $2);")
	setLockOut = ql.MustCompile("delete from Lock where Context == $1 && Task == $2;")
)

func GetLock(context Context) (task Task, ok bool, err error) {
	var lock Lock
	if err = query(nil, getLock, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&lock, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}, context.ID); err != nil {
		return
	}
	if (lock == Lock{}) {
		return
	}
	ok = true
	task, err = Get(context, lock.Task, true)
	if err != nil {
		return
	}
	return
}

func SetLockIn(context Context, task Task) (err error) {
	if context.ID != task.Context {
		err = fmt.Errorf("task's context and context doesn't match")
		return
	}
	lock := Lock{
		Context: context.ID,
		Task:    task.EID,
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, setLockIn, ql.MustMarshal(lock)...); err != nil {
		return
	}
	return
}

func SetLockOut(context Context, task Task) (err error) {
	if context.ID != task.Context {
		err = fmt.Errorf("task's context and context doesn't match")
		return
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, setLockOut, context.ID, task.EID); err != nil {
		return
	}
	return
}
