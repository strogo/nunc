package nunc

import (
	"github.com/cznic/ql"
	"strings"
	"time"
)

type Log struct {
	Command string
	Params  string
	Done    time.Time
}

var (
	trace = ql.MustCompile("insert into Log values($1, $2, $3);")
)

func Trace(command string, params []string) (err error) {
	done := time.Now()
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, trace, command, strings.Join(params, "|"), done); err != nil {
		return
	}
	return
}
