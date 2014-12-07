package nunc

import (
	"github.com/cznic/ql"
)

var (
	getLastID    = ql.MustCompile("select * from ElementID where Type == $1;")
	initEID      = ql.MustCompile("insert into ElementID values ($1, 1);")
	updateLastID = ql.MustCompile("update ElementID set Last=Last+1 where Type == $1;")
	delEID       = ql.MustCompile("delete from ElementID where Type == $1;")
)

type ElementID struct {
	Type string `ql:"index xType"`
	Last int64
}

func GetNextEID(typ string) (nextID int64, err error) {
	var eid ElementID
	if err = query(nil, getLastID, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&eid, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}, typ); err != nil {
		return
	}
	nextID = eid.Last + 1
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if (eid == ElementID{}) {
		if err = execute(ctx, initEID, typ); err != nil {
			return
		}
	} else {
		if err = execute(ctx, updateLastID, typ); err != nil {
			return
		}
	}
	return
}

func DeleteEID(typ string) (err error) {
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, delEID, typ); err != nil {
		return
	}
	return
}
