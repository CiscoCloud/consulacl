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
		"clone": func() (cli.Command, error) {
			return &command.CloneCommand{
				UI: ui,
			}, nil
		},
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
		"info": func() (cli.Command, error) {
			return &command.InfoCommand{
				UI: ui,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				UI: ui,
			}, nil
		},
		"update": func() (cli.Command, error) {
			return &command.UpdateCommand{
				UI: ui,
			}, nil
		},
	}
}
