package main

import (
	"github.com/codegangsta/cli"
	"github.com/imdario/nunc"
	"strings"
	"strconv"
)

var (
	// Command
	done = cli.Command{
		Name:      "done",
		ShortName: "x",
		Usage:     "mark a task as done",
		Action:    doneCli,
		Flags:     []cli.Flag{
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
	data := strings.SplitN(id, "-", 2)
	if len(data) != 2 {
		panic("invalid task id")
	}
	context, _, err := nunc.GetContext(data[0], true)
	if err != nil {
		panic(err)
	}
	taskId, err := strconv.ParseInt(data[1], 0, 64)
	if err != nil {
		panic(err)
	}
	beClose := c.Bool("close")
	if err := nunc.Done(context, taskId, beClose); err != nil {
		panic(err)
	}
}
