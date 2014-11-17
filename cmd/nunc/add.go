package main

import (
	"bufio"
	"github.com/codegangsta/cli"
	"github.com/imdario/nunc"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	// Command
	add = cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "add a task to context",
		Action:    addCli,
		Flags: []cli.Flag{
			text,
		},
	}
	// Add flags
	text = cli.StringFlag{
		Name: "text, t",
	}
)

const editmsgTemplate = `
# Please enter the task. Lines starting with '#' will be ignored,
# and an empty message aborts this action.`

func addCli(c *cli.Context) {
	initFromCli(c)
	defer nunc.Destroy()
	body := ""
	text := c.String("text")
	if text == "" {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			panic("Environment variable 'EDITOR' is not set")
		}
		tmp := nunc.ResolvePath("NUNC_EDITMSG")
		err := ioutil.WriteFile(tmp, []byte(editmsgTemplate), 0600)
		if err != nil {
			panic(err)
		}
		cmd := exec.Command(editor, tmp)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		file, err := os.Open(tmp)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		textCaptured := false
		lines := []string{}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "#") {
				continue
			}
			if textCaptured {
				lines = append(lines, line)
			} else {
				if line != "" {
					text = line
					textCaptured = true
				}
			}
		}
		if text == "" {
			panic("aborting new task due to empty content")
		}
		body = strings.Join(lines, "\n")
	}
	shortname := getContextFromCli(c)
	context, _, err := nunc.GetContext(shortname, true)
	if err != nil {
		panic(err)
	}
	if err = nunc.Add(context, text, body); err != nil {
		panic(err)
	}
}
