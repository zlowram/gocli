// Copyright 2014 The cli Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zlowram/gocli"
)

func testCmd1Run(tc gocli.Command, flag1 string, flag2 bool, str string) error {
	if tc.Flag.NArg() < 2 {
		tc.Usage()
		return nil
	}

	fmt.Println("This is the Run method for the subcmd1 command, with flags and args")

	fmt.Println("The value of the flag -f is:", flag1)
	if flag2 {
		fmt.Println("The bool flag -s is set")
	} else {
		fmt.Println("The bool flag -s is not set")
	}
	fmt.Println("The str value is:", str)
	fmt.Println("The first argument value is:", tc.Flag.Arg(0))
	fmt.Println("The second argument value is:", tc.Flag.Arg(1))
	return nil
}

func main() {
	flag.Parse()
	c := gocli.NewCli(
		"TestCli",
		"Short description for this cli",
		flag.Args(),
	)
	cmd1 := gocli.Command{
		Name:        "subcmd1",
		ShortName:   "s1",
		Description: "short description of the subcmd1",
		UsageLine: `TestCli subcmd1 [options] {args}

    The options are:
        -f
            string flag f default value
        -s
            bool flag s

            `,
	}
	flag1 := cmd1.Flag.String("f", "default", "")
	flag2 := cmd1.Flag.Bool("s", false, "")
	cmd1.Run = func(c gocli.Command) error {
		return testCmd1Run(c, *flag1, *flag2, "String test for Cmd1")
	}

	cmd2 := gocli.Command{
		Name:        "subcmd2",
		ShortName:   "s2",
		Description: "short description of the subcmd2",
		UsageLine: `TestCli subcmd2 [options] {args}
            `,
		Run: func(tc gocli.Command) error {
			fmt.Println("This is the Run method for the subcmd2 command, with no flags and no args")
			return nil
		},
	}
	c.Commands = []gocli.Command{cmd1, cmd2}
	if err := c.Handle(); err != nil {
		log.Fatalln(err)
	}
}
