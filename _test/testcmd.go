// Copyright 2013 The cli Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zlowram/cli"
)

type testCmd struct {
	str string
	cli.Command

	flag1 *string
}

func (tc *testCmd) Run() error {
	fmt.Println(*tc.flag1)
	fmt.Println("This is my string:", tc.str)
	return nil
}

func main() {
	flag.Parse()
	c := cli.NewCli(
		"TestCli",
		"Short description for this cli",
		flag.Args(),
	)
	cmd := testCmd{
		str: "a very long and imaginative string",
		Command: cli.Command{
			Name:        "subcmd1",
			ShortName:   "s1",
			Description: "short description of the subcmd1",
			UsageLine:   "Usage line for subcmd1"},
		//flag1: Command.Flag.String("f", "default", "string flag"),
	}
	cmd.flag1 = cmd.Flag.String("f", "default", "string flag")
	c.AddCmd(&cmd)
	c.AddCmd(&testCmd{
		str: "a very long and imaginative string 2",
		Command: cli.Command{
			Name:        "subcmd2",
			ShortName:   "s2",
			Description: "short description of the subcmd2",
			UsageLine:   "Usage line for subcmd2"},
	})
	if err := c.Handle(); err != nil {
		log.Fatalln(err)
	}
}
