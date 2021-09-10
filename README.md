# FullProxy

![build](https://img.shields.io/badge/build-passing-green)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/shoriwe/FullProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/FullProxy)](https://goreportcard.com/report/github.com/shoriwe/FullProxy)
[![Release](https://img.shields.io/github/release/shoriwe/FullProxy.svg?style=flat-square)](https://github.com/shoriwe/FullProxy/releases/latest)

\
Bind and reverse connection (with encryption) based, SOCKS5, HTTP and PortForward proxy. \

# Description

![FullProxyLogo](https://raw.githubusercontent.com/shoriwe/FullProxy/master/logo/full-proxy-logo.PNG) \
`FullProxy` is a `Bind` and `Reverse Connection` (with encryption) based `HTTP`, `SOCKS5` and `PortForward` portable
proxy

# Index

* [Title](#fullproxy)
* [Description](#description)
* [Index](#index)
* [Usage](#usage)
    * [Implemented protocols](#implemented-protocols)
        * [SOCKS5](#socks5)
        * [HTTP](#http)
        * [Forward](#port-forward)
        * [Master and Slave](#master-and-slave)
        * [Translate](#translate)
            * [Forward To SOCKS5](#port-forward-to-socks5)
    * [Implemented tools](#implemented-tools)
* [Concepts](#concepts)
    * [Master/Slave](#masterslave)
        * [How it works](#how-it-works)
        * [Applications](#applications)
    * [Translation](#translation)
* [Installation](#installation)
    * [Pre-compiled binaries](#pre-compiled-binaries)
    * [Build from source code](#build-from-source-code)
        * [Makefile](#makefile)
        * [Manual build](#manual-build)
        * [Note](#note)

# Usage

## Implemented protocols

```shell
user@linux:~$ fullproxy
```

Outputs:

```
path/to/fullproxy MODE NETWORK_TYPE ADDRESS PROTOCOL [OPTIONS]
	- MODE:         bind|master|slave
	- NETWORK_TYPE: tcp|udp
	- ADDRESS:      IPv4|IPv6 or Domain followed by ":" and the PORT; For Example -> "127.0.0.1:80"
	- PROTOCOL:     socks5|http|r-forward|l-forward|translate-socks5
Environment Variables:
	- C2Address     Host and port of the C2 port of the master server
```

### SOCKS5

```shell
user@linux:~$ fullproxy bind tcp 127.0.0.1:9050 socks5 --help
```

Outputs:

```shell
Usage of socks5:
  -auth-cmd string
        shell command to pass the hex encoded username and password, exit code 0 means login success
  -inbound-blacklist string
        plain text file list with all the HOST that are forbidden to connect to the proxy
  -inbound-whitelist string
        plain text file list with all the HOST that are permitted to connect to the proxy
  -outbound-blacklist string
        plain text file list with all the forbidden proxy targets
  -outbound-whitelist string
        plain text file list with all the permitted proxy targets
  -users-file string
        json file with username as keys and sha3-513 of the password as values
```

### HTTP

HTTP proxy could be implemented thanks to [GoProxy](https://github.com/elazarl/goproxy)

```shell
user@linux:~$ fullproxy local-forward -help
```

Outputs:

```shell
```

### Port Forward

```shell
user@linux:~$ fullproxy bind tcp 127.0.0.1:9050 port-forward --help
```

Outputs

```shell
Usage of port-forward:
  -inbound-blacklist string
        plain text file list with all the HOST that are forbidden to connect to the proxy
  -inbound-whitelist string
        plain text file list with all the HOST that are permitted to connect to the proxy
  -network-type string
        tcp or udp (default "tcp")
  -target-address string
        Address to connect (default "127.0.0.1:80")
```

### Master and Slave

#### Preparing the master

Serve a socks5 proxy using master and slave.

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && fullproxy master tcp 127.0.0.1:9050 socks5 [OPTIONS]
```

#### Preparing the slave

Connect to the master and proxy the networks from the slave side.

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && fullproxy slave tcp 127.0.0.1:9050 socks5
```

### Translate

#### Port Forward To SOCKS5

```shell
user@linux:~$ fullproxy translate port_forward-socks5 -help
```

Outputs:

```shell
```

## Implemented tools

# Concepts

## Master/Slave

Handles the proxying between a reverse connected proxy and the clients. In other words, it will receive the connections
of the clients and will forward the traffic to the proxy that is reverse connected to it.

### How it works

1. It first binds to the host specified by the user.
2. Then accept the connection from the proxy server.
3. Finally, it proxy the traffic of all new incoming connections to the proxy server that was reverse connected to it in
   the second step. In other words, is the proxy of another proxy but totally invisible for the client.

### Applications

This could be specially useful when you need to proxy a network that a machine have access to, but you can't bind inside
that machine.

## Translation

This protocol is simple, it receives proxying request in a specific proxying protocol to them forward them to another
proxy with another protocol; this means that if you only speaks SOCKS5, you will be able to talk to an HTTP proxy using
this "translator"

# Installation

## Pre-compiled binaries

You can find pre-compiled binaries for windows and linux [Here](https://github.com/shoriwe/FullProxy/releases)

## Build from source code

```shell
go install github.com/shoriwe/FullProxy/cmd/fullproxy@latest
```

### Note

For some reason (possibly a coding fault commit by me) if you extremely stress the master/slave protocol, it will crash
in away that it is still running and the new connections are received but the connections to the targets are never made.