# FullProxy

![build](https://img.shields.io/badge/build-passing-green)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/shoriwe/FullProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/FullProxy)](https://goreportcard.com/report/github.com/shoriwe/FullProxy)
[![Release](https://img.shields.io/github/release/shoriwe/FullProxy.svg?style=flat-square)](https://github.com/shoriwe/FullProxy/releases/latest)

Bind and reverse connection based, SOCKS5, HTTP and PortForward proxy.

# Description

![logo](logo/white_logo_color_background.jpg)

`FullProxy` is a `Bind` and `Reverse Connection` (with TLS) `HTTP`, `SOCKS5` and `PortForward` portable proxy.

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
    * [Implemented tools](#implemented-tools)
* [Concepts](#concepts)
    * [Master/Slave](#masterslave)
        * [How it works](#how-it-works)
        * [Applications](#applications)
* [Installation](#installation)
    * [Pre-compiled binaries](#pre-compiled-binaries)
    * [Build from source code](#build-from-source-code)
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
 - PROTOCOL:     socks5|http|port-forward|translate-socks5
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

HTTP proxy could be implemented thanks to [goproxy](https://github.com/elazarl/goproxy)

```shell
user@linux:~$ fullproxy http -help
```

Outputs:

```shell
Usage of http:
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

Serve a socks5 proxy using master and slave.

Notice that if you don't specify a certificate (with environment variable `C2Certificate=/path/to/cert`) and a private
key (with environment variable `C2PrivateKey=/path/to/priv.key`) in the master command, the tool will automatically
generate one for you.

If you are using a self-signed certificate, or you did not specify one, you should also use the environment
variable `C2SlaveIgnoreTrust=1` to continue on untrusted cert.

#### Setup without certificate

##### Preparing the master

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && fullproxy master tcp 127.0.0.1:9050 socks5 [OPTIONS]
```

##### Preparing the slave

Connect to the master and proxy the networks from the slave side.

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && export C2SlaveIgnoreTrust="1" && fullproxy slave tcp 127.0.0.1:9050 socks5
```

#### Setup with certificate

##### Preparing the master

```shell
user@linux:~$ export C2PrivateKey=/path/to/priv.key && export C2Certificate=/path/to/cert && export C2Address="127.0.0.1:9051" && fullproxy master tcp 127.0.0.1:9050 socks5 [OPTIONS]
```

##### Preparing the slave

Connect to the master and proxy the networks from the slave side.

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && fullproxy slave tcp 127.0.0.1:9050 socks5
```

Notice that if your certificate is still invalid for the client you should try:

```shell
user@linux:~$ export C2Address="127.0.0.1:9051" && export C2SlaveIgnoreTrust="1" && fullproxy slave tcp 127.0.0.1:9050 socks5
```

## Implemented tools

### fullproxy-users

This tool will create a valid `JSON` file to use with the flag `-users-file`

```shell
user@linux:~$ fullproxy-users
```

Outputs:

```shell
fullproxy-users COMMAND DATABASE_FILE USERNAME
Available commands:
        - new
        - delete
        - set
```

#### Commands:

##### new

Creates a new file with a new user.

##### delete

Deletes an existing user in the file.

##### set

Creates or updates a user in the file.

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

# Installation

## Pre-compiled binaries

You can find pre-compiled binaries for windows and Linux [Here](https://github.com/shoriwe/FullProxy/releases)

## Build from source code

### fullproxy

```shell
go install github.com/shoriwe/fullproxy/v3/cmd/fullproxy@latest
```

### fullproxy-users

```shell
go install github.com/shoriwe/fullproxy/v3/cmd/fullproxy-users@latest
```

### Note

For some reason (possibly a coding fault commit by me) if you extremely stress the master/slave protocol, it will crash
in away that it is still running and the new connections are received but the connections to the targets are never made.