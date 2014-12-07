package main

import (
	"fmt"
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	get = cli.Command{
		Name:      "get",
		ShortName: "g",
		Usage:     "get a task from context",
		Action:    getCli,
		Flags: []cli.Flag{
			any,
			nolock,
		},
	}
	// Add flags
	any = cli.BoolFlag{
		Name: "any, a",
	}
	nolock = cli.BoolFlag{
		Name: "nolock, n",
	}
)

func getCli(c *cli.Context) {
	context, taskId, err := getTaskAndContextFromCli(c)
	if err != nil {
		panic(err)
	}
	var (
		ok   bool
		task nunc.Task
	)
	lock := !c.Bool("nolock")
	if lock {
		task, ok, err = nunc.GetLock(context)
		if err != nil {
			panic(err)
		}
	}
	if !ok {
		if taskId == 0 {
			if c.Bool("any") {
				task, context, err = nunc.GetOldestFromAllContexts()
			} else {
				task, err = nunc.GetOldest(context)
			}
		} else {
			task, err = nunc.Get(context, taskId, true)
		}
		if err != nil {
			panic(err)
		}
		if lock {
			err = nunc.SetLockIn(context, task)
			if err != nil {
				panic(err)
			}
		}
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
