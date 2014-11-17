package nunc

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
)

type State int32

const (
	Closed State = iota - 1
	Resolved
	Open
	InProgress
	Feedback
)

type Task struct {
	ID           int64 `ql:"index xID"`
	Context      int64
	Text         string
	State        State
	Creation     time.Time
	Modification time.Time
	Parent       int64
}

func (t *Task) StateString() (state string) {
	switch t.State {
	case Closed:
		state = "x"
	case Resolved:
		state = "."
	case Open:
		state = "o"
	case InProgress:
		state = ">"
	case Feedback:
		state = "?"
	default:
		state = "!"
	}
	return
}

func TaskID(context Context, task Task) string {
	return fmt.Sprintf("@%s-%d", context.ShortName, task.ID)
}

func TaskPath(context Context, task Task) string {
	return filepath.Join(context.ShortName, strconv.FormatInt(task.ID, 10))
}

func TaskBody(context Context, task Task) (string, error) {
	data, err := ioutil.ReadFile(ResolvePath(TaskPath(context, task)))
	return string(data), err
}
