package main

import (
	"github.com/codegangsta/cli"
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
	initFromCli(c)
	defer nunc.Destroy()
	id := c.Args().First()
	context, taskId, err := nunc.ResolveTaskID(id)
	if err != nil {
		panic(err)
	}
	beClose := c.Bool("close")
	if err := nunc.Done(context, taskId, beClose); err != nil {
		panic(err)
	}
}
