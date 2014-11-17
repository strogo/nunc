package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/imdario/nunc"
	"strings"
	"strconv"
)

var (
	// Command
	get = cli.Command{
		Name:      "get",
		ShortName: "g",
		Usage:     "get a task from context",
		Action:    getCli,
	}
)

func getCli(c *cli.Context) {
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
	task, err := nunc.Get(context, taskId, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s]\n\n%s\n\n", nunc.TaskID(context, task), task.Text)
	body, err := nunc.TaskBody(context, task)
	if err != nil {
		panic(err)
	}
	if body != "" {
		fmt.Println(body)
	}
}
