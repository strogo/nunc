package main

import (
	"github.com/codegangsta/cli"
)

var (
	// Command
	edit = cli.Command{
		Name:      "edit",
		ShortName: "e",
		Usage:     "edit a task from context",
		Action:    editCli,
	}
)

func editCli(c *cli.Context) {

}
