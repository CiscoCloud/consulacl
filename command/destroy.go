package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type DestroyCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *DestroyCommand) Help() string {
	helpText := `
Usage: consulacl destroy [options] id

  Destroy an ACL

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

func (c *DestroyCommand) Run(args []string) int {
	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	extra := cmdFlags.Args()
	if len(extra) < 1 {
		c.UI.Error("ACL id must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	id := extra[0]

	consul, err := NewConsulClient(c.Consul, &c.UI)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.ACL()

	writeOpts := new(consulapi.WriteOptions)

	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *DestroyCommand) Synopsis() string {
	return "Destroy an ACL"
}
