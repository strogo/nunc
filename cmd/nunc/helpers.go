package main

import (
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
	"os"
	"strings"
)

func getTaskAndContextFromCli(c *cli.Context) (context nunc.Context, taskId int64, err error) {
	id := c.Args().First()
	context, taskId, err = nunc.ResolveTaskID(id)
	if err != nil {
		return
	}
	if context.ShortName == "" {
		shortname := resolveContext("")
		context, _, err = nunc.GetContext(shortname, true)
		if err != nil {
			return
		}
	}
	return
}

func getContextFromCli(c *cli.Context) (shortname string) {
	shortname = c.Args().First()
	return resolveContext(shortname)
}

func resolveContext(raw string) (shortname string) {
	shortname = raw
	if raw == "" {
		shortname = os.Getenv("NUNC_MAIN_CONTEXT")
	}
	shortname = strings.Map(dropSlash, shortname)
	return
}

func dropSlash(r rune) rune {
	if r == '-' {
		return -1
	}
	return r
}
