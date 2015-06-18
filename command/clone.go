package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type CloneCommand struct {
	UI	cli.Ui
	Consul	*ConsulFlags
}

func (c *CloneCommand) Help() string {
	helpText := `
Usage: consulacl clone [options] id

  Create a new token from an existing one

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

func (c *CloneCommand) Run(args []string) int {
	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	extra := cmdFlags.Args()
	if len(extra) < 1 {
		c.UI.Error("ACL id must be provided")
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

	newid, _, err := client.Clone(id, nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(newid)

	return 0
}

func (c *CloneCommand) Synopsis() string {
	return "Create a new token from an existing one"
}
