package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type CreateCommand struct {
	UI		cli.Ui
	Consul		*ConsulFlags
	ConfigRules	[]*ConfigRule
}

type ConfigRule struct {
	PathType	string
	Path		string
	Policy		string
}

func (c *CreateCommand) Help() string {
	helpText := `
Usage: consulacl create [options]

  Create an ACL. Requires a management token.

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
  --management			Create a management token
				(default: false)
  --name			Name of the ACL
				(default: not set)
  --rule='type:path:policy'	Rule to create. Can be multiple rules on a command line
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *CreateCommand) Run(args []string) int {
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

//	if len(c.ConfigRules) < 1 && !isManagement {
//		c.UI.Error("Must supply an acl rule")
//		c.UI.Error("")
//		c.UI.Error(c.Help())
//		return 1
//	}

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
			Name:	aclName,
			Type:	consulapi.ACLClientType,
			Rules:	rules,
		}

	}

	id, _, err := client.Create(entry, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(id)

	return 0
}

func (c *CreateCommand) Synopsis() string {
	return "Create an ACL"
}

func ParseRuleConfig(s string) (*ConfigRule, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, errors.New("cannot specify empty rule declaration")
	}

	var pathType, path, policy string
	parts := strings.Split(s, ":")

	switch len(parts) {
	case 2:
		pathType, path = parts[0], parts[1]
		policy = "read"
	case 3:
		pathType, path, policy = parts[0], parts[1], parts[2]
	default:
		return nil, fmt.Errorf("invalid rule declaration '%s'", s)
	}

	return &ConfigRule{ pathType, path, policy }, nil
}

type rulePath struct {
	Policy	string
}

type aclRule struct {
	Key	map[string]*rulePath	`json:"key,omitempty"`
	Service	map[string]*rulePath	`json:"service,omitempty"`
}

// Convert a list of Rules to a JSON string
func GetRulesString(rs []*ConfigRule)  (string, error) {
	rules := &aclRule{
		Key:		make(map[string]*rulePath),
		Service:	make(map[string]*rulePath),
	}

	for _, r := range rs {
		// Verify policy is one of "read", "write", or "deny"
		policy := strings.ToLower(r.Policy)
		switch policy {
		case "read", "write", "deny":
		default:
			return "", fmt.Errorf("Invalid rule policy: '%s'", r.Policy)
		}

		switch strings.ToLower(r.PathType) {
		case "key":
			rules.Key[r.Path] = &rulePath{ r.Policy }
		case "service":
			rules.Service[r.Path] = &rulePath{ r.Policy }
		default:
			return "", fmt.Errorf("Invalid path type: '%s'", r.PathType)
		}
	}

	ruleBytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}

	return string(ruleBytes), nil
}
