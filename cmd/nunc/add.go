package main

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/imdario/cli"
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
			paste,
		},
	}
	// Add flags
	text = cli.StringFlag{
		Name: "text, t",
	}
	paste = cli.BoolFlag{
		Name: "clipboard, c",
	}
)

const editmsgTemplate = `
# Please enter the task. Lines starting with '#' will be ignored,
# and an empty message aborts this action.`

func addCli(c *cli.Context) {
	shortname := getContextFromCli(c)
	context, _, err := nunc.GetContext(shortname, true)
	if err != nil {
		panic(err)
	}
	body, text := getTextAndBodyFromContext(c)
	task, err := nunc.Add(context, text, body)
	if err != nil {
		panic(err)
	}
	fmt.Println("New task:", nunc.TaskID(context, task))
}

func getTextAndBodyFromContext(c *cli.Context) (body, text string) {
	text = c.String("text")
	if text == "" {
		var (
			tmp string
			err error
		)
		if c.Bool("clipboard") {
			tmp, err = openClipboard()
		} else {
			tmp, err = openEditor()
		}
		if err != nil {
			panic(err)
		}
		text, body, err = readTaskFrom(tmp)
		if err != nil {
			panic(err)
		}
	}
	return
}

func openClipboard() (tmp string, err error) {
	tmp = nunc.ResolvePath("NUNC_EDITMSG")
	text, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		err = fmt.Errorf("aborting new task due to empty content")
		return
	}
	fmt.Println("Do you want to add this as task? (y/N)")
	elipsis := ""
	length := len(text)
	if length > 140 {
		length = 140
		elipsis = "..."
	}
	fmt.Printf("\n%s%s\n\n> ", text[:length], elipsis)
	reader := bufio.NewReader(os.Stdin)
	input, _, _ := reader.ReadRune()
	switch input {
	case 'y', 'Y':
		// noop
	default:
		err = fmt.Errorf("aborting new task due to user input")
		return
	}
	err = ioutil.WriteFile(tmp, []byte(text), 0600)
	if err != nil {
		return
	}
	return
}

func openEditor() (tmp string, err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		panic("Environment variable 'EDITOR' is not set")
	}
	tmp = nunc.ResolvePath("NUNC_EDITMSG")
	err = ioutil.WriteFile(tmp, []byte(editmsgTemplate), 0600)
	if err != nil {
		return
	}
	cmd := exec.Command(editor, tmp)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	return
}

func readTaskFrom(tmp string) (text, body string, err error) {
	file, err := os.Open(tmp)
	if err != nil {
		return
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
			if len(lines) == 0 {
				if line == "" {
					continue
				}
			}
			lines = append(lines, line)
		} else {
			if line != "" {
				text = line
				textCaptured = true
			}
		}
	}
	if text == "" {
		err = fmt.Errorf("aborting new task due to empty content")
		return
	}
	body = strings.Join(lines, "\n")
	return
}
