package nunc

import (
	"fmt"
	"github.com/cznic/ql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type internalState struct {
	Home        string
	Initialized bool
	Db          *ql.DB
}

func ResolvePath(file string) (path string) {
	return filepath.Join(is.Home, file)
}

func ResolveTaskID(id string) (context Context, taskId int64, err error) {
	data := strings.SplitN(id, "-", 2)
	if len(data) != 2 {
		err = fmt.Errorf("invalid task id")
	}
	context, _, err = GetContext(data[0], true)
	if err != nil {
		return
	}
	taskId, err = strconv.ParseInt(data[1], 0, 64)
	if err != nil {
		return
	}
	return
}

var is internalState

func Init(home string) (err error) {
	if is.Initialized {
		return
	}
	is.Home = home
	if err = os.MkdirAll(is.Home, 0750); err != nil {
		return
	}
	if err = initDb(); err != nil {
		return
	}
	is.Initialized = true
	return
}

func initDb() (err error) {
	opts := ql.Options{
		CanCreate: true,
	}
	is.Db, err = ql.OpenFile(ResolvePath("nunc.db"), &opts)
	if err != nil {
		return
	}
	ctx := ql.NewRWCtx()
	initDbTable(ctx, (*Context)(nil))
	if err != nil {
		return
	}
	initDbTable(ctx, (*Tag)(nil))
	if err != nil {
		return
	}
	initDbTable(ctx, (*Task)(nil))
	if err != nil {
		return
	}
	initDbTable(ctx, (*Tagged)(nil))
	if err != nil {
		return
	}
	return
}

func initDbTable(ctx *ql.TCtx, specimen interface{}) (err error) {
	schema := ql.MustSchema(specimen, "", nil)
	_, _, err = is.Db.Execute(ctx, schema)
	return
}

func Destroy() (err error) {
	err = is.Db.Flush()
	if err != nil {
		return
	}
	err = is.Db.Close()
	if err != nil {
		return
	}
	is.Initialized = false
	return
}
