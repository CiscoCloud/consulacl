package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

const Name = "consulacl"
const Version = "0.1.1"

func main() {
	log.SetOutput(ioutil.Discard)

	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "--version" {
			fmt.Printf("%s v%s\n", Name, Version)
			os.Exit(0)
		}
	}

	cli := &cli.CLI{
		Args:		args,
		Commands:	Commands,
		HelpFunc:	cli.BasicHelpFunc("consulacl"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(exitCode)
}
