// Copyright 2014 The cli Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zlowram/cli"
)

type testCmd1 struct {
	str string
	cli.Command

	flag1 *string
	flag2 *bool
}

func (tc *testCmd1) Run() error {
	fmt.Println("This is the Run method for the subcmd1 command, with flags and args")

	fmt.Println("The value of the flag -f is:", *tc.flag1)
	if *tc.flag2 {
		fmt.Println("The bool flag -s is set")
	} else {
		fmt.Println("The bool flag -s is not set")
	}
	fmt.Println("The str value is:", tc.str)
	return nil
}

type testCmd2 struct {
	cli.Command
}

func (tc *testCmd2) Run() error {
	fmt.Println("This is the Run method for the subcmd2 command, with no flags and no args")
	return nil
}

func main() {
	flag.Parse()
	c := cli.NewCli(
		"TestCli",
		"Short description for this cli",
		flag.Args(),
	)
	cmd1 := testCmd1{
		str: "a very long and imaginative string",
		Command: cli.Command{
			Name:        "subcmd1",
			ShortName:   "s1",
			Description: "short description of the subcmd1",
			UsageLine: `TestCli subcmd1 [options] {args}

    The options are:
        -f
            string flag f default value
        -s
            bool flag s

            `},
	}
	cmd1.flag1 = cmd1.Flag.String("f", "default", "")
	cmd1.flag2 = cmd1.Flag.Bool("s", false, "")
	cmd2 := testCmd2{
		Command: cli.Command{
			Name:        "subcmd2",
			ShortName:   "s2",
			Description: "short description of the subcmd2",
			UsageLine: `TestCli subcmd2 [options] {args}
            `},
	}
	c.AddCmds([]cli.CliCommand{
		&cmd1,
		&cmd2,
	})
	if err := c.Handle(); err != nil {
		log.Fatalln(err)
	}
}
