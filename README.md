# consulacl
Command line interface to the [Consul ACL HTTP API](https://consul.io/docs/agent/http/acl.html). Documentation for the Consul ACL system is at [the Consul ACL internals][Consul ACLs] page.

## Installation
You can download a released `consulacl` artifact from [the consulacl release page][Releases] on Github. If you wish to compile from source, you will need to have buildtools and [Go][] installed:

```shell
$ git clone https://github.com/CiscoCloud/consulacl.git
$ cd consulacl
$ make
```

## Basic Usage

```shell
usage: consulacl [--version] [--help] <command> [<args>]

Available commands are:
    clone      Create a new token from an existing one
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

#### Example

```shell
$ consulacl clone --sll --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 19933651-439e-5123-5a2f-6bdf2afa0d70
a06db641-070d-eae0-1ff8-8e8c67399fa4
```

### create command

#### Usage

```shell
Usage: consulacl create [options]

  Create an ACL. Requires a management token.

Options:

  --management			Create a management token
				(default: false)
  --name			Name of the ACL
				(default: not set)
  --rule='type:path:policy'	Rule to create. Can be multiple rules on a command line
				(default: not set)
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `management` | `false` | Create the token as a management ACL
| `name` | `not set` | Name of the ACL
| `rule` | `not set` | Rule to create

Multiple rules can be specified on the command line.  The format for the `rule` is `[key|service]:path:[read:write:deny]`. The list of rules is converted to a JSON object:

```JSON
{
  "key": {
    "<path_1>": {
      "policy": "<policy_1>"
    }, ...
   },
  "service": {
    "<path_2>": {
      "policy": "<policy_2>"
    }, ...
  }
}
```

An empty `path` attribute generates:

```JSON
{
   "key": {
     "": {
       "policy": "<policy_1>"
     }
   }
}
```

The token id of the newly created ACL is printed on stdout on success.

#### Example

```shell
$ consulacl create --ssl --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 \
    --rule='key:test/node:read' \
    --rule='service:hello-world:write'
25c25096-e680-2faa-d864-b9314308387a
```

### destroy command

#### Usage

```shell
consulacl destroy [options] id

  Destroy an ACL
```

#### Example

```shell
$ consulacl destroy --ssl --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 \
    25c25096-e680-2faa-d864-b9314308387a
```

### info command

#### Usage

```shell
consulacl info [options] id

  Query information about an ACL token
```

#### Example

```shell
$ consulacl info --ssl --ssl-verify=false --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 \
	 25c25096-e680-2faa-d864-b9314308387a
{
  "CreateIndex": 4100,
  "ModifyIndex": 4100,
  "ID": "25c25096-e680-2faa-d864-b9314308387a",
  "Name": "",
  "Type": "client",
  "Rules": "{\"key\":{\"test/node\":{\"Policy\":\"read\"}},\"service\":{\"hello-world\":{\"Policy\":\"write\"}}}"
}

### list command

#### Usage

```shell
consulacl list [options]

  List all active ACL tokens.
```

#### Example

```shell
$ consulacl list --ssl --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 
{
  {
    "CreateIndex": 3,
    "ModifyIndex": 3,
    "ID": "anonymous",
    "Name": "Anonymous Token",
    "Type": "client",
    "Rules": ""
  },
  {
    "CreateIndex": 4100,
    "ModifyIndex": 4100,
    "ID": "25c25096-e680-2faa-d864-b9314308387a",
    "Name": "",
    "Type": "client",
    "Rules": "{\"key\":{\"test/node\":{\"Policy\":\"read\"}},\"service\":{\"hello-world\":{\"Policy\":\"write\"}}}"
  }
}
```

### update command

The update command updates an ACL if it exists and creates a new one if it does not. All of the ACL settings are overwritten on update.

#### Usage

```shell
Usage: consulacl update [options] id

  Update an ACL. Will be created if it doesn't exist.

Options:

  --management			Create a management token
				(default: false)
  --name			Name of the ACL
				(default: not set)
  --rule='type:path:policy'	Rule to create. Can be multiple rules on a command line
				(default: not set)
```

#### Arguments

| Option | Default | Description |
| ------ | ------- | ----------- |
| `management` | `false` | Create the token as a management ACL
| `name` | `not set` | Name of the ACL
| `rule` | `not set` | Rule to create

Multiple rules can be specified on the command line.  The format for the `rule` is `[key|service]:path:[read:write:deny]`. The list of rules is converted to a JSON object:

```JSON
{
  "key": {
    "<path_1>": {
      "policy": "<policy_1>"
    }, ...
   },
  "service": {
    "<path_2>": {
      "policy": "<policy_2>"
    }, ...
  }
}
```

An empty `path` attribute generates:

```JSON
{
   "key": {
     "": {
       "policy": "<policy_1>"
     }
   }
}
```

The token id of the newly created ACL is printed on stdout on success.

#### Example

```shell
$ consulacl update --ssl --token=b78191f9-01fb-24d0-4278-be05ee82c6c4 \
    --rule='key:test/node:read' \
    --rule=`key:test/node1:write' \
    --rule='service:hello-world:write' \
    25c25096-e680-2faa-d864-b9314308387a
```

[Consul ACLs]: http://www.consul.io/docs/internals/acl.html "Consul ACLs"
[Releases]: https://github.com/CiscoCloud/consulacl/releases "consulacl releases page"
[Go]: http://golang.org "Go the language"
