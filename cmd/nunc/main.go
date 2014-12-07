package main

import (
	"github.com/imdario/cli"
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
		done,
		fork,
		context,
	}
	app.Before = initFromCli
	app.After = destroyFromCli
	app.Flags = []cli.Flag{
		home,
	}
	app.Run(os.Args)
}

func initFromCli(c *cli.Context) (err error) {
	home := c.GlobalString("home")
	err = nunc.Init(home)
	if err != nil {
		return
	}
	return
}

func destroyFromCli(c *cli.Context) (err error) {
	args := c.Args()
	nunc.Trace(args[0], args[1:])
	nunc.Destroy()
	return
}
