package main

import (
	"os"
	"github.com/codegangsta/cli"
)

var (
	get = cli.Command{
		Name: "get",
		ShortName: "g",
		Usage: "get a task from context",
		Action: getCli,
	}
	list = cli.Command{
		Name: "list",
		ShortName: "l",
		Usage: "list contexts and their tasks",
		Action: listCli,
	}
	add = cli.Command{
		Name: "add",
		ShortName: "a",
		Usage: "add a task to context",
		Action: addCli,
	}
	del = cli.Command{
		Name: "delete",
		ShortName: "d",
		Usage: "delete a task from context",
		Action: deleteCli,
	}
	edit = cli.Command{
		Name: "edit",
		ShortName: "e",
		Usage: "edit a task from context",
		Action: editCli,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "nunc"
	app.Usage = "no more procastination"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		get,
		list,
		add,
		del,
		edit,
	}
	app.Run(os.Args)
}
