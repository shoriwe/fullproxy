# FullProxy
 [![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/FullProxy)](https://goreportcard.com/report/github.com/shoriwe/FullProxy)
 [![Generic badge](https://img.shields.io/badge/Releases-ALL-any.svg)](https://github.com/shoriwe/FullProxy/releases)

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
* [Concepts](#concepts)
    * [Master/Slave](#masterslave)
        * [How it works](#how-it-works)
        * [Applications](#applications)
        * [Considerations](#considerations)
    * [Translation](#translation)
* [Installation](#installation)
    * [Pre-compiled binaries](#pre-compiled-binaries)
    * [Build from source code](#build-from-source-code)
* [Suggestions](#suggestions)
# Usage
## Implemented protocols
```
user@linux:~$ fullproxy help
Usage:
         fullproxy.exe PROTOCOL *FLAGS

Protocols available:
         - socks5
         - http
         - local-forward
         - remote-forward
         - master
         - translate
```
### SOCKS5
```
user@linux:~$ fullproxy socks5 --help
Usage of socks5:
  -host string
        Host to listen on. When "-slave" flag is set, is the IP of master to connect
  -password string
        Password of the running proxy, requires "-username". It will be ignored if is an empty string
  -port string
        Port to listen on. When "-slave" flag is set, is the Port of the master to connect. I both modes the default port is 1080 (default "1080")
  -slave
        Connect to a master, no bind proxying
  -username string
        Username of the running proxy, requires "-password". It will be ignored if is an empty string
```
### HTTP
HTTP proxy could be implemented thanks to [GoProxy](https://github.com/elazarl/goproxy)
```
user@linux:~$ fullproxy local-forward -help
Usage of http:
  -host string
        Host to listen on. When "-slave" flag is set, is the IP of master to connect
  -password string
        Password of the running proxy, requires "-username". It will be ignored if is an empty string
  -port string
        Port to listen on. When "-slave" flag is set, is the Port of the master to connect. I both modes the default port is 8080 (default "8080")
  -slave
        Connect to a master, no bind proxying
  -tls
        Use HTTPS
  -username string
        Username of the running proxy, requires "-password". It will be ignored if is an empty string
```
### Forward
#### Local
```
user@linux:~$ fullproxy local-forward -help
Usage of local-forward:
  -forward-host string
        Host to forward the traffic received from master
  -forward-port string
        Port to forward the traffic received from master
  -master-host string
        Host of the master
  -master-port string
        Port of the master
```
#### Remote
```
user@linux:~$ fullproxy remote-forward -help
Usage of remote-forward:
  -local-host string
        Host to bind by slave
  -local-port string
        Port to bind by slave
  -master-host string
        Host of the master
  -master-port string
        Port of the master
```
### Master
```
user@linux:~$ fullproxy remote-forward -help
Usage of master:
  -host string
        Host to listen on. (default "0.0.0.0")
  -forward-host string
        Argument required to handle correctly the "remote-forward" (This is the service that the master can only acceded)
  -forward-port string
        Argument required to handle correctly the "remote-forward" (This is the service that the master can only acceded)
  -port string
        Port to listen on. (default "1080")
```
### Translate
```
user@linux:~$ fullproxy translate help
Usage:
         fullproxy.exe translate TARGET *FLAGS

TARGETS available:
         - port_forward-socks5
```
#### Port Forward To SOCKS5
```
user@linux:~$ fullproxy translate port_forward-socks5 -help
Usage of port_forward-socks5:
  -bind-host string
        Host to listen on. (default "0.0.0.0")
  -bind-port string
        Port to listen on. (default "8080")
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
### Considerations
- The `master` protocol may loss some setup connections if it is extremely stressed, but it should work `just fine` if the connections where already made
## Translation
This protocol is simple, it receives proxying request in a specific proxying protocol to them forward them to another proxy with another protocol; this means that if you only speaks SOCKS5, you will be able to talk to an HTTP proxy using this "translator" 
# Installation
## Pre-compiled binaries
You can find pre-compiled binaries for windows and linux [Here](https://github.com/shoriwe/FullProxy/releases)
## Build from source code
Run this command:
```
go get github.com/shoriwe/FullProxy
```
# Suggestions
If you have any suggestion for new features, also leave them in the issue section or create the proper branch, add what do you want and request a pull request
