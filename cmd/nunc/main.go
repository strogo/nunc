package main

import (
	"github.com/codegangsta/cli"
	"github.com/imdario/nunc"
	"log"
	"os"
)

var (
	// Logging
	logger = log.New(os.Stderr, "nunc: ", 0)
)

func main() {
	if os.Getenv("NUNC_DEV") == "" {
		defer func() {
			if r := recover(); r != nil {
				logger.Fatal(r)
			}
		}()
	}
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
		context,
		done,
	}
	app.Flags = []cli.Flag{
		home,
	}
	app.Run(os.Args)
}

func initFromCli(c *cli.Context) {
	home := c.GlobalString("home")
	if err := nunc.Init(home); err != nil {
		panic(err)
	}
}
