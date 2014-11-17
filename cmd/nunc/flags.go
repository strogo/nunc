package main

import (
	"github.com/codegangsta/cli"
	"os/user"
	"path/filepath"
)

var (
	// Global flags
	home = cli.StringFlag{
		Name:   "home, m",
		Value:  defaultHomeFlag(),
		EnvVar: "NUNC_HOME",
	}
	verbose = cli.BoolFlag{
		Name: "verbose, v",
	}
	all = cli.BoolFlag{
		Name: "all, a",
	}
)

func defaultHomeFlag() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, ".nunc")
}
