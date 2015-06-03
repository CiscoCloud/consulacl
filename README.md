# consulacl
Command line interface to the [Consul ACL HTTP API](https://consul.io/docs/agent/http/acl.html). Documentation for the Consul ACL system is at [the Consul ACL internals][Consul ACLs] page.

## Installation
You can download a released `consulacl` artifact from [the consulacl release page][Releases] on Github. If you wish to compile from source, you will need to have buildtools and [Go][] installed:

```shell
$ git clone https://github.com/CiscoCloud/consulkv.git
$ cd consulacl
$ make
```

## Basic Usage

```shell
usage: consulacl [--version] [--help] <command> [<args>]

Available commands are:
    [clone](#clone command)      Create a new token from an existing one
    create     Create an ACL
    destroy    Destroy an ACL
    info       Query an ACL token
    list       List a value
    update     Update an ACL
```

### Common arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `--consul` | `127.0.0.1:8500` | HTTP address of the Consul Agent
| `--ssl` | `false` | Use HTTPS while talking to Consul
| `--ssl-verify` | `true` | Verify certificates when connecting via SSL. Requires `--ssl`
| `--ssl-cert` | `unset` | Path to an SSL client certificate to use to authenticate to the consul server
| `--ssl-ca-cert` | `unset` | Path to a CA certificate file, containing one or more CA certificates to use to validate the certificate sent by the consul server to us.
| `--token`* | `unset` | The [Consul API token][Consul ACLs].

\* A management token is required for all ACL operations

### clone command

#### Usage

```shell
consulacl clone [options] id

  Create a new token from an existing one
```


[Consul ACLs]: http://www.consul.io/docs/internals/acl.html "Consul ACLs"
[Releases]: https://github.com/CiscoCloud/consulacl/releases "consulacl releases page"
[Go]: http://golang.org "Go the language"
