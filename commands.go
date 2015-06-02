package main

import (
	"os"

	"github.com/CiscoCloud/consulacl/command"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{ Writer: os.Stdout }

	Commands = map[string]cli.CommandFactory{
		"create": func() (cli.Command, error) {
			return &command.CreateCommand{
				UI: ui,
			}, nil
		},
		"destroy": func() (cli.Command, error) {
			return &command.DestroyCommand{
				UI: ui,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				UI: ui,
			}, nil
		},
	}
}
