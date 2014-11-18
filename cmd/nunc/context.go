package main

import (
	"fmt"
	"github.com/imdario/cli"
	"github.com/imdario/nunc"
)

var (
	// Command
	context = cli.Command{
		Name:      "context",
		ShortName: "c",
		Usage:     "options for contexts",
		Subcommands: []cli.Command{
			contextList,
			contextAdd,
			contextDelete,
			contextEdit,
			contextPurge,
		},
	}
	// Subcommands
	contextList = cli.Command{
		Name:      "list",
		ShortName: "l",
		Usage:     "list existing contexts",
		Action:    contextListCli,
		Flags: []cli.Flag{
			verbose,
			all,
		},
	}
	contextAdd = cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "add a new context",
		Action:    contextAddCli,
		Flags: []cli.Flag{
			contextName,
		},
	}
	contextDelete = cli.Command{
		Name:      "delete",
		ShortName: "d",
		Usage:     "delete a context",
		Action:    contextDeleteCli,
	}
	contextEdit = cli.Command{
		Name:      "edit",
		ShortName: "e",
		Usage:     "edit a context",
		Action:    contextEditCli,
		Flags: []cli.Flag{
			contextName,
			contextShortName,
		},
	}
	contextPurge = cli.Command{
		Name:      "purge",
		ShortName: "p",
		Usage:     "purge an inactive context",
		Action:    contextPurgeCli,
	}
	// Add flags
	contextName = cli.StringFlag{
		Name:  "name, n",
		Usage: "long name for the context",
	}
	// Edit flags
	contextShortName = cli.StringFlag{
		Name:  "shortname, s",
		Usage: "short name for the context",
	}
)

func contextListCli(c *cli.Context) {
	beAll := c.Bool("all")
	contexts, err := nunc.ListContexts(beAll)
	beVerbose := c.Bool("verbose")
	if err != nil {
		panic(err)
	}
	for _, context := range contexts {
		data := []interface{}{context.ShortName}
		if beVerbose {
			data = append(data, "("+context.Name+")")
		}
		if beAll {
			if context.Inactive {
				data = append(data, "[v]")
			}
		}
		fmt.Println(data...)
	}
}

func contextAddCli(c *cli.Context) {
	shortname := getContextFromCli(c)
	name := c.String("name")
	if err := nunc.AddContext(name, shortname); err != nil {
		panic(err)
	}
}

func contextDeleteCli(c *cli.Context) {
	shortname := getContextFromCli(c)
	if err := nunc.DeleteContext(shortname); err != nil {
		panic(err)
	}
}

func contextEditCli(c *cli.Context) {
	oldShortname := getContextFromCli(c)
	name := c.String("name")
	shortname := c.String("shortname")
	if name == "" && shortname == "" {
		panic("you must provide a -n or -s flag")
	}
	if err := nunc.EditContext(name, oldShortname, shortname); err != nil {
		panic(err)
	}
}

func contextPurgeCli(c *cli.Context) {
	shortname := getContextFromCli(c)
	if err := nunc.PurgeContext(shortname); err != nil {
		panic(err)
	}
}
