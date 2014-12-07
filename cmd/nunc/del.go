package main

import (
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	del = cli.Command{
		Name:      "del",
		ShortName: "d",
		Usage:     "delete a task from context",
		Action:    delCli,
	}
)

func delCli(c *cli.Context) {
	context, taskId, err := getTaskAndContextFromCli(c)
	if err != nil {
		panic(err)
	}
	_, err = nunc.Delete(context, taskId)
	if err != nil {
		panic(err)
	}
}
