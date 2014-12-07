package main

import (
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	fork = cli.Command{
		Name:      "fork",
		ShortName: "f",
		Usage:     "fork a task from context",
		Action:    forkCli,
		Flags: []cli.Flag{
			text,
		},
	}
)

func forkCli(c *cli.Context) {
	context, taskId, err := getTaskAndContextFromCli(c)
	if err != nil {
		panic(err)
	}
	task, ok, err := nunc.GetLock(context)
	if err != nil {
		panic(err)
	}
	if ok && task.Parent > 0 {
		panic("you cannot fork tasks created as forks")
	}
	body, text := getTextAndBodyFromContext(c)
	newTask, err := nunc.Fork(context, text, body, taskId)
	if err != nil {
		panic(err)
	}
	if ok {
		err = nunc.SetLockOut(context, task)
		if err != nil {
			panic(err)
		}
		err = nunc.SetLockIn(context, newTask)
		if err != nil {
			panic(err)
		}
	}
}
