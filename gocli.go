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

// The CliCommand interface defines the methods that a Command must implement. 
type CliCommand interface {
	GetName() string
	GetShortName() string
	GetDescription() string
	GetUsageLine() string
	getFlag() *flag.FlagSet

	Usage()
	Run() error
}

// Cli represents the whole CLI program, which includes the name, description,
// commands and arguments.
type Cli struct {
	commands    []CliCommand
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

// AddCmds add commands to the current Cli. Note that all commands must
// implement the CliCommand interface.
func (cl *Cli) AddCmds(cmds []CliCommand) {
	for _, cmd := range cmds {
		cl.commands = append(cl.commands, cmd)
	}
}

// Usage prints the help text for the current Cli.
func (cl *Cli) Usage() {
	clTemplate := struct {
		Commands    []CliCommand
		Name        string
		Description string
	}{cl.commands, cl.name, cl.description}

	if err := usageTemplate.Execute(os.Stdout, clTemplate); err != nil {
		log.Fatalln(err)
	}
}

var usageTemplate = template.Must(template.New("usage").Parse(usageLine))

const usageLine = `{{.Name}} - {{.Description}}

    Usage:
        command [options] {arguments}

    The commands are:
        {{range .Commands}} {{.GetName}}, {{.GetShortName}}     {{.GetDescription}}
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
	for _, cmd := range cl.commands {

		f := cmd.getFlag()
		f.Usage = func() { cmd.Usage(); os.Exit(1) }
		f.Parse(cl.args[1:])

		switch cl.args[0] {
		case cmd.GetName():
			fallthrough
		case cmd.GetShortName():
			return cmd.Run()
		case "help":
			if cmd.GetName() == cl.args[1] || cmd.GetShortName() == cl.args[1] {
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
}

// GetName returns the name of the command.
func (c Command) GetName() string {
	return c.Name
}

// GetName returns the short name of the command.
func (c Command) GetShortName() string {
	return c.ShortName
}

// GetDescription returns the description of the command.
func (c Command) GetDescription() string {
	return c.Description
}

// GetUsageLine returns the usage text of the command.
func (c Command) GetUsageLine() string {
	return c.UsageLine
}

func (c *Command) getFlag() *flag.FlagSet {
	return &c.Flag
}

// Usage prints the usage text of the command to the stdout.
func (c Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.GetUsageLine())
}
