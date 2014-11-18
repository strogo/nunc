package nunc

import (
	"time"
)

type Log struct {
	ID      int64 `ql:"index xID"`
	Command string
	Context int64
	Task    int64
	Params  string
	Done    time.Time
}
