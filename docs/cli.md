# CLI

This document present a description of the features provided by the command line interface of `fullproxy`.

## HELP

To print the help message of CLI as first argument is necessary the word `help`.

```shell
fullproxy help
```

Expected output:

```
Usage: fullproxy COMMAND [ARGUMENTS]

Available commands:

- help:			prints this help.
- slave:		Connects to master server.
- socks5:		Starts a SOCKS5 server.
- http:			Starts a HTTP proxy server.
- forward:		Starts a port forward proxy server.
- translate:	Translate a proxy protocol to another to proxy protocol.
- reverse:		Starts a raw reverse proxy.
- config:		Start serving the server configured in the targeted yaml file.
```

To obtain the help of any command, you should execute it without arguments or with `-h`, `--help` arguments.

## Slave

The slave command connects to master to forward traffic.

```shell
Usage of translate:
  -master string
        Address of master server. Argument URL structure is 'network://host:port'
```

## SOCKS5

The `socks5` command permits the user to start quickly a SOCKS5 server.

Commands help looks like:

```shell
fullproxy socks5 --help
```

Expected output:

```
Usage of socks5:
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

## HTTP

The `http` command permits the users to start quickly an HTTP proxy server.

Command help looks like:

```shell
fullproxy http --help
```

Expected output:

```
Usage of http:
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

## Port Forward

The `forward` command permits the users to start quickly a port forward proxy server.

Command help looks like:

```shell
fullproxy forward --help
```

Expected output:

```
Usage of forward:
  -dial-tls
        Dial connection will use TLS
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
  -target string
        Target forward address. Argument URL structure is 'network://host:port'
```

## Translate

The `translate` command permits the users to start quickly a translation proxy server.

Command help looks like:

```shell
fullproxy translate --help
```

Expected output:

```
Usage of translate:
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
  -source string
        Address of the source proxy. Argument URL structure is 'network://host:port/protocol'
  -target string
        Address of target proxy. Argument URL structure is 'network://host:port/protocol'
```

## Reverse

The `reverse` command permits the users to start quickly a raw port reverse server.

Command help looks like:

```shell
fullproxy reverse --help
```

Expected output:

```
Usage of reverse:
  -add value
        Add target to the pool load balancer
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

## Config

The `config` command starts all the proxy server specified in a YAML file. See `yaml.md` for more information.

Command help looks like:

```shell
fullproxy config
```

Expected output:

```
Usage: fullproxy config YAML_CONFIG
```

