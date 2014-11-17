package nunc

import (
	"github.com/cznic/ql"
)

var (
	beginStmt    = ql.MustCompile("begin transaction;")
	commitStmt   = ql.MustCompile("commit;")
	rollbackStmt = ql.MustCompile("rollback;")
)

type dbCallback func(data []interface{}) (bool, error)

func beginTransaction() (ctx *ql.TCtx, err error) {
	ctx = ql.NewRWCtx()
	_, _, err = is.Db.Execute(ctx, beginStmt)
	if err != nil {
		return
	}
	return
}

func commit(ctx *ql.TCtx) (err error) {
	_, _, err = is.Db.Execute(ctx, commitStmt)
	if err != nil {
		return
	}
	return
}

func rollback(ctx *ql.TCtx) (err error) {
	_, _, err = is.Db.Execute(ctx, rollbackStmt)
	if err != nil {
		return
	}
	return
}

func query(ctx *ql.TCtx, stmt ql.List, cb dbCallback, arg ...interface{}) (err error) {
	rs, _, err := is.Db.Execute(ctx, stmt, arg...)
	if err != nil {
		return
	}
	if err = rs[0].Do(false, cb); err != nil {
		return
	}
	return
}

func execute(ctx *ql.TCtx, stmt ql.List, arg ...interface{}) (err error) {
	_, _, err = is.Db.Execute(ctx, stmt, arg...)
	return
}

func autoclose(ctx *ql.TCtx, err error) (nerr error) {
	if err == nil {
		nerr = commit(ctx)
	} else {
		rollback(ctx)
	}
	return
}
