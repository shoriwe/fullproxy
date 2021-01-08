# FullProxy
![build](https://img.shields.io/badge/build-passing-green)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/shoriwe/FullProxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/FullProxy)](https://goreportcard.com/report/github.com/shoriwe/FullProxy)
[![Release](https://img.shields.io/github/release/shoriwe/FullProxy.svg?style=flat-square)](https://github.com/shoriwe/FullProxy/releases/latest)

 \
 Bind and reverse connection (with encryption) based, SOCKS5, HTTP and PortForward proxy.
 \
# Description
![FullProxyLogo](https://raw.githubusercontent.com/shoriwe/FullProxy/master/logo/full-proxy-logo.PNG) \
`FullProxy` is a `Bind` and `Reverse Connection` (with encryption) based `HTTP`, `SOCKS5` and `PortForward` portable proxy
# Index
* [Title](#fullproxy)
* [Description](#description)
* [Index](#index)
* [Usage](#usage)
    * [Implemented protocols](#implemented-protocols)
        * [SOCKS5](#socks5)
        * [HTTP](#http)
        * [Forward](#forward)
        * [Master](#master)
        * [Translate](#translate)
            * [Forward To SOCKS5](#port-forward-to-socks5)
    * [Implemented tools](#implemented-tools)
        * [Database](#database)
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
* [Suggestions](#suggestions)
# Usage
## Implemented protocols
```shell
user@linux:~$ fullproxy help
Usage:
         fullproxy PROTOCOL|TOOL *FLAGS

Protocols available:
         - socks5
         - http
         - local-forward
         - remote-forward
         - master
         - translate

Tools available:
         - database
```
### SOCKS5
```shell
user@linux:~$ fullproxy socks5 --help
Usage of socks5:
  -command-auth string
        Command with it's default args to pass the Username and Password received from clients, please notice that ExitCode = 0 will mean that the login was successful, any other way i
t not and the username and password will be passed as base64 encoded arguments to it, this auth method will ignore any other supplied
  -database-auth string
        Path to the SQLite3 database generated with the 'database create' command and filled with the 'database user add' command, this auth method will ignore any other supplied
  -host string
        Host to listen on. When "-slave" flag is set, is the IP of master to connect
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)
  -outbound-blacklist string
        File with a host per line. Denied outgoing connections (ignored when outbound-whitelist is set)
  -outbound-whitelist string
        File with a host per line. Allowed outgoing connections (ignored when outbound-blacklist is set)
  -password-auth string
        Password of the running proxy, requires "-username". It will be ignored if is an empty string, this auth method will ignore any other supplied
  -port string
        Port to listen on. When "-slave" flag is set, is the Port of the master to connect. I both modes the default port is 1080 (default "1080")
  -slave
        Connect to a master, no bind proxying
  -timeout int
        The number of second before re-trying the connection between target and client (default is 10 seconds) (default 10)
  -tries int
        The number of re-tries that will maintain the connection between target and client (default is 5 tries) (default 5)
  -username-auth string
        Username of the running proxy, requires "-password". It will be ignored if is an empty string, this auth method will ignore any other supplied
```
### HTTP
HTTP proxy could be implemented thanks to [GoProxy](https://github.com/elazarl/goproxy)
```shell
user@linux:~$ fullproxy local-forward -help
Usage of http:
  -command-auth string
        Command with it's default args to pass the Username and Password received from clients, please notice that ExitCode = 0 will mean that the login was successful, any other way i
t not and the username and password will be passed as base64 encoded arguments to it, this auth method will ignore any other supplied
  -database-auth string
        Path to the SQLite3 database generated with the 'database create' command and filled with the 'database user add' command, this auth method will ignore any other supplied
  -host string
        Host to listen on. When "-slave" flag is set, is the IP of master to connect
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)
  -outbound-blacklist string
        File with a host per line. Denied outgoing connections (ignored when outbound-whitelist is set)
  -outbound-whitelist string
        File with a host per line. Allowed outgoing connections (ignored when outbound-blacklist is set)
  -password-auth string
        Password of the running proxy, requires "-username". It will be ignored if is an empty string, this auth method will ignore any other supplied
  -port string
        Port to listen on. When "-slave" flag is set, is the Port of the master to connect. I both modes the default port is 8080 (default "8080")
  -slave
        Connect to a master, no bind proxying
  -tls
        Use HTTPS
  -username-auth string
        Username of the running proxy, requires "-password". It will be ignored if is an empty string, this auth method will ignore any other supplied
```
### Forward
#### Local
```shell
user@linux:~$ fullproxy local-forward -help
Usage of local-forward:
  -forward-host string
        Host to forward the traffic received from master
  -forward-port string
        Port to forward the traffic received from master
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored when inbound-blacklist is set)
  -master-host string
        Host of the master
  -master-port string
        Port of the master
  -timeout int
        The number of second before re-trying the connection between target and client (default is 10 seconds) (default 10)
  -tries int
        The number of re-tries that will maintain the connection between target and client (default is 5 tries) (default 5)
```
#### Remote
```shell
user@linux:~$ fullproxy remote-forward -help
Usage of remote-forward:
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored when inbound-blacklist is set)
  -local-host string
        Host to bind by slave
  -local-port string
        Port to bind by slave
  -master-host string
        Host of the master
  -master-port string
        Port of the master
  -timeout int
        The number of second before re-trying the connection between target and client (default is 10 seconds) (default 10)
  -tries int
        The number of re-tries that will maintain the connection between target and client (default is 5 tries) (default 5)
```
### Master
```shell
user@linux:~$ fullproxy remote-forward -help
Usage of master:
  -forward-host string
        Argument required to handle correctly the "remote-forward" (This is the service that the master can only acceded)
  -forward-port string
        Argument required to handle correctly the "remote-forward" (This is the service that the master can only acceded)
  -host string
        Host to listen on. (default "0.0.0.0")
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored when inbound-blacklist is set)
  -port string
        Port to listen on. (default "1080")
  -timeout int
        The number of second before re-trying the connection between target and client (default is 10 seconds) (default 10)
  -tries int
        The number of re-tries that will maintain the connection between target and client (default is 5 tries) (default 5)
```
### Translate
```shell
user@linux:~$ fullproxy translate help
Usage:
         fullproxy translate TARGET *FLAGS

TARGETS available:
         - port_forward-socks5
```
#### Port Forward To SOCKS5
```shell
user@linux:~$ fullproxy translate port_forward-socks5 -help
Usage of port_forward-socks5:
  -bind-host string
        Host to listen on. (default "0.0.0.0")
  -bind-port string
        Port to listen on. (default "8080")
  -inbound-blacklist string
        File with a host per line. Denied incoming connections to the proxy (ignored when inbound-whitelist is set)
  -inbound-whitelist string
        File with a host per line. Allowed incoming connections to the proxy (ignored when inbound-blacklist is set)
  -socks5-host string
        SOCKS5 server host to use (default "127.0.0.1")
  -socks5-password string
        Password for the SOCKS5 server; leave empty for no AUTH
  -socks5-port string
        SOCKS5 server port to use (default "1080")
  -socks5-username string
        Username for the SOCKS5 server; leave empty for no AUTH
  -target-host string
        Host of the target host that is accessible by the SOCKS5 proxy
  -target-port string
        Port of the target host that is accessible by the SOCKS5 proxy
  -timeout int
        The number of second before re-trying the connection between target and client (default is 10 seconds) (default 10)
  -tries int
        The number of re-tries that will maintain the connection between target and client (default is 5 tries) (default 5)
```
## Implemented tools
### Database
This tool helps the user in the creation and administration of `SQLite3` database with the actual structure that `FullProxy` supports
```shell
user@linux:~$ fullproxy database help
Usage:
         fullproxy database CMD

CMDs available:
         - create
         - user
```
#### Create
Tools here are used to maintain an already created database
```shell
user@linux:~$ fullproxy database create
Usage:
        fullproxy database create DATABASE_FILE
```
#### User
```shell
user@linux:~$ fullproxy database user help
Usage:
         fullproxy database user CMD

CMDs available:
         - add
         - update
         - delete
```
##### Add
```shell
user@linux:~$ fullproxy database user add help
Usage:
        fullproxy database user add DATABASE_FILE USERNAME PASSWORD
```
##### Delete
```shell
user@linux:~$ fullproxy database user delete help
Usage:
        fullproxy database user delete DATABASE_FILE USERNAME
```
##### Update
```shell
user@linux:~$ fullproxy database user update help
Usage:
        fullproxy database user update DATABASE_FILE USERNAME NEW_PASSWORD
```
# Concepts
## Master/Slave
Handles the proxying between a reverse connected (with encryption) proxy and the clients. In other words, it will receive the connections of the clients and will forward the traffic to the proxy that is reverse connected to it.
### How it works
1. It first binds to the host specified by the user.
2. Then accept the connection from the proxy server.
3. Finally, it proxy the traffic of all new incoming connections to the proxy server that was reverse connected to it in the second step.
In other words, is the proxy of another proxy but totally invisible for the client.
### Applications
This could be specially useful when you need to proxy a network that a machine have access to, but you can't bind with it
## Translation
This protocol is simple, it receives proxying request in a specific proxying protocol to them forward them to another proxy with another protocol; this means that if you only speaks SOCKS5, you will be able to talk to an HTTP proxy using this "translator" 
# Installation
## Pre-compiled binaries
You can find pre-compiled binaries for windows and linux [Here](https://github.com/shoriwe/FullProxy/releases)
## Build from source code
### Makefile
You can approach the `Makefile` that I prepare for the project, you just need to set the environment variables `CC` and `CXX` and compiled based on:
```bash
make OS-ARCH-LINKING
```
For example:
- Compiling a static binary for a 64-bit based linux
```bash
make linux-64-static
```
- Compiling a dynamic binary for 32-bit based windows
```bash
make windows-32-dynamic
```
### Manual build
- Download the source code:
```shell
go get github.com/shoriwe/FullProxy
```
- Go to `cmd/FullProxy`
```shell
cd ~/go/src/github.com/shoriwe/FullProxy/cmd/FullProxy
```
- Compile it
```shell
# Statically
CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -linkmode external -extldflags=-static" -tags sqlite_omit_load_extension -mod vendor
# Or Dynamically
CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -mod vendor
```
### Note
In some systems it will be better to dynamically compile the binary instead of statically and in others, the other way, this probably happens because in how each manage it's networking features and/or the dependencies of the sqlite3 library
# Suggestions
If you have any suggestion for new features, also leave them in the issue section or create the proper branch, add what do you want and request a pull request
