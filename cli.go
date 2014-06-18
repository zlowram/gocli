package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
)

type CliCommand interface {
	GetName() string
	GetShortName() string
	GetDescription() string
	GetUsageLine() string
	GetFlag() flag.FlagSet

	Usage()
	Run() error
}

type Cli struct {
	commands    []CliCommand
	args        []string
	name        string
	description string
}

func NewCli(name string, description string, args []string) *Cli {
	return &Cli{
		name:        name,
		description: description,
		args:        args,
	}
}

func (cl *Cli) AddCmd(cmd CliCommand) {
	cl.commands = append(cl.commands, cmd)
}

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

func (cl *Cli) Handle() error {
	if len(cl.args) < 1 || len(cl.args) == 1 && cl.args[0] == "help" {
		cl.Usage()
		return nil
	}
	for _, cmd := range cl.commands {

        f := cmd.GetFlag()
        f.Usage = func() { f.PrintDefaults(); os.Exit(1) }
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

type Command struct {
	Name        string
	ShortName   string
	Description string
	UsageLine   string
	Flag        flag.FlagSet
}

func (c Command) GetName() string {
	return c.Name
}

func (c Command) GetShortName() string {
	return c.ShortName
}

func (c Command) GetDescription() string {
	return c.Description
}

func (c Command) GetUsageLine() string {
	return c.UsageLine
}

func (c Command) GetFlag() flag.FlagSet {
	return c.Flag
}

func (c Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", c.GetUsageLine())
}
