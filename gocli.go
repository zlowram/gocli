// Copyright 2014 The cli Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

// Cli represents the whole CLI program, which includes the name, description,
// commands and arguments.
type Cli struct {
	Commands    []Command
	args        []string
	name        string
	description string
}

// NewCli creates and returns a new Cli. In order to create a new CLI program,
// this function should be called in the first place.
func NewCli(name string, description string, args []string) *Cli {
	return &Cli{
		name:        name,
		description: description,
		args:        args,
	}
}

// Usage prints the help text for the current Cli.
func (cl *Cli) Usage() {
	clTemplate := struct {
		Commands    []Command
		Name        string
		Description string
	}{cl.Commands, cl.name, cl.description}

	if err := usageTemplate.Execute(os.Stdout, clTemplate); err != nil {
		log.Fatalln(err)
	}
}

var usageTemplate = template.Must(template.New("usage").Parse(usageLine))

const usageLine = `{{.Name}} - {{.Description}}

    Usage:
        command [options] {arguments}

    The commands are:
        {{range .Commands}} {{.Name}}, {{.ShortName}}     {{.Description}}
        {{end}}
    Use "{{.Name}} help [command]" for more information about a command.

`

// Handle is the engine of CLI. It handles the arguments given to the program
// and calls the corresponding command Run or Usage methods.
func (cl *Cli) Handle() error {
	if len(cl.args) < 1 || len(cl.args) == 1 && cl.args[0] == "help" {
		cl.Usage()
		return nil
	}
	for _, cmd := range cl.Commands {

		if cl.args[0] != cmd.Name && cl.args[0] != cmd.ShortName && cl.args[0] != "help" {
			continue
		}

		cmd.Flag.Usage = func() { cmd.Usage(); os.Exit(1) }
		cmd.Flag.Parse(cl.args[1:])

		switch cl.args[0] {
		case cmd.Name:
			fallthrough
		case cmd.ShortName:
			return cmd.Run(cmd)
		case "help":
			if cmd.Name == cl.args[1] || cmd.ShortName == cl.args[1] {
				cmd.Usage()
				return nil
			}
		}
	}
	if err := unknownCmdTemplate.Execute(os.Stdout, cl.name); err != nil {
		log.Fatalln(err)
	}
	return nil
}

var unknownCmdTemplate = template.Must(template.New("unknown").Parse(unknownCmdLine))

const unknownCmdLine = `unknown command

Use "{{.}} help" for usage information.

`

// Command represents a command of the cli program. It includes the name,
// short name, description, usage and flags.
type Command struct {
	Name        string
	ShortName   string
	Description string
	UsageLine   string
	Flag        flag.FlagSet
	Run         func(Command) error
}

// Usage prints the usage text of the command to the stdout.
func (c Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.UsageLine)
}
