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
Usage: fullproxy slave MASTER_ADDRESS
```

## SOCKS5

The `socks5` command permits the user to start quickly a SOCKS5 server.

Commands help looks like:

```shell
fullproxy socks5 --help
```

Expected output:

```
Usage: fullproxy socks5 [ARGUMENTS]

Arguments:

-h, --help: 	Show this help message.
--listen="":	Address to listen.
--master="":	Listen address for master/slave communication.
```

## HTTP

The `http` command permits the users to start quickly an HTTP proxy server.

Command help looks like:

```shell
fullproxy http --help
```

Expected output:

```
Usage: fullproxy http [ARGUMENTS]

Arguments:

-h, --help: 	Show this help message.
--listen="":	Address to listen.
--master="":	Listen address for master/slave communication.
```

## Port Forward

The `forward` command permits the users to start quickly a port forward proxy server.

Command help looks like:

```shell
fullproxy forward --help
```

Expected output:

```
Usage: fullproxy forward [ARGUMENTS]

Arguments:

-h, --help: 	Show this help message.
--listen="":	Address to listen.
--master="":	Listen address for master/slave communication.
--target="":	Target address to redirect the traffic.
```

## Translate

The `translate` command permits the users to start quickly a translation proxy server.

Command help looks like:

```shell
fullproxy translate --help
```

Expected output:

```
Usage: fullproxy translate [ARGUMENTS]

Arguments:

-h, --help: 			Show this help message.
--listen="":			Address to listen.
--source=""				Address to the source proxy.
--source-protocol="":	Proxy protocol to used by clients.
--target="":			Target proxy address.
--target-protocol="":	Proxy protocol been used by target.
--master="":			Listen address for master/slave communication.
```

## Reverse

The `reverse` command permits the users to start quickly a raw port reverse server.

Command help looks like:

```shell
fullproxy reverse --help
```

Expected output:

```
Usage: fullproxy reverse [ARGUMENTS]

Arguments:

-h, --help: 	Show this help message.
--pool:			List of targets used by the load balancer.
--listen="":	Address to listen.
--master="":	Listen address for master/slave communication.
```

## Config

The `config` command starts all the proxy server specified in a YAML file. See `yaml.md` for more information.

Command help looks like:

```shell
fullproxy config --help
```

Expected output:

```
Usage: fullproxy config PATH_TO_YAML_FILE
```

