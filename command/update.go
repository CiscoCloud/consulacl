package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type UpdateCommand struct {
	UI		cli.Ui
	Consul		*ConsulFlags
	ConfigRules	[]*ConfigRule
}

func (c *UpdateCommand) Help() string {
	helpText := `
Usage: consulacl update [options] id

  Update an ACL. Will be created if it doesn't exist.

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
  --management			Update type to 'management'
				(default: false)
  --name			Name of the ACL
				(default: not set)
  --rule='type:path:policy'	Rule to create. Can be multiple rules on a command line
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *UpdateCommand) Run(args []string) int {
	var isManagement bool
	var aclName string

	c.Consul = new(ConsulFlags)
	cmdFlags := NewFlagSet(c.Consul)
	cmdFlags.StringVar(&aclName, "name", "", "")
	cmdFlags.BoolVar(&isManagement, "management", false, "")

	cmdFlags.Var((funcVar)(func(s string) error {
		t, err := ParseRuleConfig(s)
		if err != nil {
			return err
		}

		if c.ConfigRules == nil {
			c.ConfigRules = make([]*ConfigRule, 0, 1)
		}

		c.ConfigRules = append(c.ConfigRules, t)
		return nil
	}), "rule", "")

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

	var entry *consulapi.ACLEntry

	if isManagement {
		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	aclName,
			Type:	consulapi.ACLManagementType,
		}
	} else {
		rules, err := GetRulesString(c.ConfigRules)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	aclName,
			Type:	consulapi.ACLClientType,
			Rules:	rules,
		}

	}

	_, err = client.Update(entry, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *UpdateCommand) Synopsis() string {
	return "Update an ACL"
}
