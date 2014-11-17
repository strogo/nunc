package main

import (
	"github.com/codegangsta/cli"
)

var (
	// Command
	del = cli.Command{
		Name:      "delete",
		ShortName: "d",
		Usage:     "delete a task from context",
		Action:    deleteCli,
	}
)

func deleteCli(c *cli.Context) {

}
