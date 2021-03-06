package main

import (
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	done = cli.Command{
		Name:      "done",
		ShortName: "x",
		Usage:     "mark a task as done",
		Action:    doneCli,
		Flags: []cli.Flag{
			closed,
		},
	}
	// Done flags
	closed = cli.BoolFlag{
		Name: "close, c",
	}
)

func doneCli(c *cli.Context) {
	context, taskId, err := getTaskAndContextFromCli(c)
	if err != nil {
		panic(err)
	}
	beClose := c.Bool("close")
	if err := nunc.Done(context, taskId, beClose); err != nil {
		panic(err)
	}
	task, err := nunc.Get(context, taskId, true)
	if err != nil {
		panic(err)
	}
	err = nunc.SetLockOut(context, task)
	if err != nil {
		panic(err)
	}
}
