package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	list = cli.Command{
		Name:      "list",
		ShortName: "l",
		Usage:     "list contexts and their tasks",
		Action:    listCli,
		Flags: []cli.Flag{
			all,
		},
	}
)

func listCli(c *cli.Context) {
	initFromCli(c)
	defer nunc.Destroy()
	shortname := getContextFromCli(c)
	context, _, err := nunc.GetContext(shortname, true)
	if err != nil {
		panic(err)
	}
	beAll := c.Bool("all")
	tasks, err := nunc.ListTasks(context, beAll)
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		state := task.StateString()
		fmt.Println(
			state,
			nunc.TaskID(context, task),
			"::",
			task.Text,
		)
	}
}
