package command

import (
	"encoding/json"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type ListCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *ListCommand) Help() string {
	helpText := `
Usage: consulacl list [options]

  List all active ACL tokens.

Options:

  --consul=127.0.0.1:8500	HTTP address of the Consul Agent
  --ssl				Use HTTPS while talking to Consul.
				(default: false)
  --ssl-verify			Verify certificates when connecting via SSL.
				(default: true)
  --ssl-cert			Path to an SSL certificate to use to authenticate
				to the consul server.
				(default: not set)
  --ssl-ca-cert			Path to an SSL client certificate to use to authenticate
				to the consul server.
				(default: not set)
  --token			The Consul API token.
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Run(args []string) int {
	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	queryOpts := new(consulapi.QueryOptions)
	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.ACL()

	acls, _, err := client.List(queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	jsonRaw, err := json.MarshalIndent(acls, "", "  ")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(string(jsonRaw))

	return 0
}

func (c *ListCommand) Synopsis() string {
	return "List a value"
}
