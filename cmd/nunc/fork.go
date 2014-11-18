package main

import (
	"github.com/imdario/cli"
)

var (
	// Command
	fork = cli.Command{
		Name:      "fork",
		ShortName: "f",
		Usage:     "fork a task from context",
		Action:    forkCli,
	}
)

func forkCli(c *cli.Context) {

}
