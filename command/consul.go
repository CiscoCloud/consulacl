package command

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	flag "github.com/ogier/pflag"
	"github.com/mitchellh/cli"
)

type ConsulFlags struct {
	consulAddr	string
	sslEnabled	bool
	sslVerify	bool
	sslCert		string
	sslCaCert	string
	token		string
	auth		*Auth
}

func NewFlagSet(c *ConsulFlags) *flag.FlagSet {
	consulFlags := flag.NewFlagSet("consulkv", flag.ContinueOnError)
	consulFlags.StringVar(&c.consulAddr, "consul", "127.0.0.1:8500", "")
	consulFlags.BoolVar(&c.sslEnabled, "ssl", false, "")
	consulFlags.BoolVar(&c.sslVerify, "ssl-verify", true, "")
	consulFlags.StringVar(&c.sslCert, "ssl-cert", "", "")
	consulFlags.StringVar(&c.sslCaCert, "ssl-ca-cert", "", "")
	consulFlags.StringVar(&c.token, "token", "", "")

	c.auth = new(Auth)
	consulFlags.Var((*Auth)(c.auth), "auth", "")

	return consulFlags
}

func NewConsulClient(c *ConsulFlags, ui *cli.Ui) (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddr

	if c.token != "" {
		config.Token = c.token
	}

	if c.sslEnabled {
		config.Scheme = "https"
	}

	if !c.sslVerify {
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if c.auth.Enabled {
		config.HttpAuth = &consulapi.HttpBasicAuth{
			Username: c.auth.Username,
			Password: c.auth.Password,
		}
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authentication var
type Auth struct {
	Enabled		bool
	Username	string
	Password	string
}


func (a *Auth) Set(value string) error {
	a.Enabled = true

	if (strings.Contains(value, ":")) {
		split := strings.SplitN(value, ":", 2)
		a.Username = split[0]
		a.Password = split[1]
	} else {
		a.Username = value
	}

	return nil
}

func (a *Auth) String() string {
	if a.Password == "" {
		return a.Username
	}

	return fmt.Sprintf("%s:%s", a.Username, a.Password)
}

type funcVar func(s string) error

func (f funcVar) Set(s string) error	{ return f(s) }
func (f funcVar) String() string	{ return "" }
func (f funcVar) IsBoolFlag() bool	{ return false }
