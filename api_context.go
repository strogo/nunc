package nunc

import (
	"errors"
	"fmt"
	"github.com/cznic/ql"
	"os"
	"strings"
)

var (
	getContext           = ql.MustCompile("select * from Context where ShortName == $1;")
	listContexts         = ql.MustCompile("select * from Context where Inactive == FALSE;")
	listAllContexts      = ql.MustCompile("select * from Context;")
	addContext           = ql.MustCompile("insert into Context values($1, $2, $3);")
	delContext           = ql.MustCompile("update Context set Inactive = TRUE where ShortName == $1;")
	editContextName      = ql.MustCompile("update Context set Name = $1 where ShortName == $2;")
	editContextShortName = ql.MustCompile("update Context set ShortName = $1 where ShortName == $2;")
	purgeContextTasks    = ql.MustCompile("delete from Task where Context == $1;")
	purgeContext         = ql.MustCompile("delete from Context where ShortName == $1;")
)

func sanitizeContext(shortname string) (clean string, err error) {
	clean = strings.TrimSpace(shortname)
	if clean == "" {
		err = errors.New("context not provided")
		return
	}
	if !strings.HasPrefix(clean, "@") {
		err = errors.New("invalid context, must be prefixed by '@'")
		return
	}
	clean = clean[1:]
	if clean == "" {
		err = errors.New("context was provided but empty")
		return
	}
	return
}

func GetContext(shortname string, must bool) (context Context, id string, err error) {
	id, err = sanitizeContext(shortname)
	if err != nil {
		return
	}
	if err = query(nil, getContext, func(data []interface{}) (bool, error) {
		if er2 := ql.Unmarshal(&context, data); er2 != nil {
			return false, er2
		}
		return true, nil
	}, id); err != nil {
		return
	}
	if (context == Context{}) && must {
		err = fmt.Errorf("unknown context '@%s'", id)
	}
	return
}

func ListContexts(all bool) (contexts []Context, err error) {
	q := listContexts
	if all {
		q = listAllContexts
	}
	if err = query(nil, q, func(data []interface{}) (bool, error) {
		context := Context{}
		if er2 := ql.Unmarshal(&context, data); er2 != nil {
			return false, er2
		}
		contexts = append(contexts, context)
		return true, nil
	}); err != nil {
		return
	}
	return
}

func AddContext(name, shortname string) (err error) {
	context, id, err := GetContext(shortname, false)
	if err != nil {
		return
	}
	if (context != Context{}) {
		err = fmt.Errorf("context '%s' already exists", id)
		return
	}
	context = Context{
		Name:      name,
		ShortName: id,
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, addContext, ql.MustMarshal(context)...); err != nil {
		return
	}
	if err = os.MkdirAll(ResolvePath(id), 0750); err != nil {
		return
	}
	return
}

func DeleteContext(shortname string) (err error) {
	_, id, err := GetContext(shortname, true)
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
	if err = execute(ctx, delContext, id); err != nil {
		return
	}
	return
}

func EditContext(name, oldShortname, shortname string) (err error) {
	_, id, err := GetContext(oldShortname, true)
	if err != nil {
		return
	}
	context, newId, err := GetContext(shortname, false)
	if err != nil {
		return
	}
	if (context != Context{}) {
		err = fmt.Errorf("context '%s' already exists", newId)
		return
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if name != "" {
		err = updateContextName(ctx, name, id)
		if err != nil {
			return
		}
	}
	if shortname != "" {
		err = updateContextShortName(ctx, newId, id)
		if err != nil {
			return
		}
	}
	return
}

func updateContextName(ctx *ql.TCtx, name, id string) (err error) {
	if err = execute(ctx, editContextName, name, id); err != nil {
		return
	}
	return
}

func updateContextShortName(ctx *ql.TCtx, newId, id string) (err error) {
	if err = execute(ctx, editContextShortName, newId, id); err != nil {
		return
	}
	if err = os.Rename(ResolvePath(id), ResolvePath(newId)); err != nil {
		return
	}
	return
}

func PurgeContext(shortname string) (err error) {
	context, id, err := GetContext(shortname, true)
	if err != nil {
		return
	}
	if !context.Inactive {
		err = fmt.Errorf("context '%s' must be inactive to be purged", id)
		return
	}
	ctx, err := beginTransaction()
	if err != nil {
		return
	}
	defer func() {
		autoclose(ctx, err)
	}()
	if err = execute(ctx, purgeContextTasks, context.ID); err != nil {
		return
	}
	if err = execute(ctx, purgeContext, context.ShortName); err != nil {
		return
	}
	if err = os.RemoveAll(ResolvePath(id)); err != nil {
		return
	}
	return
}
