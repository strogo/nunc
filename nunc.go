package nunc

import (
	"github.com/cznic/ql"
	"os"
)

type internalState struct {
	Home        string
	Initialized bool
	Db          *ql.DB
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
	initDbTable(ctx, (*Lock)(nil))
	if err != nil {
		return
	}
	initDbTable(ctx, (*Log)(nil))
	if err != nil {
		return
	}
	initDbTable(ctx, (*ElementID)(nil))
	if err != nil {
		return
	}
	// TODO meta
	return
}

func initDbTable(ctx *ql.TCtx, specimen interface{}) (err error) {
	schema := ql.MustSchema(specimen, "", nil)
	_, _, err = is.Db.Execute(ctx, schema)
	return
}

func Destroy() (err error) {
	if !is.Initialized {
		return
	}
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
