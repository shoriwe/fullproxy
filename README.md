# FullProxy
 Bind and reverse connection (with encryption) based, SOCKS5, HTTP and PortForward proxy.
 
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
* [Concepts](#concepts)
    * [Slave](#slave)
        * [How it works](#how-it-works)
        * [Applications](#applications)
        * [Considerations](#considerations)
* [Installation](#installation)
    * [Pre-compiled binaries](#pre-compiled-binaries)
    * [Build from source code](#build-from-source-code)
* [Suggestions](#suggestions)
# Usage
## Implemented protocols
```
user@linux:~$ fullproxy help
Usage:  fullproxy  PROTOCOL *FLAGS
Protocols available:
         - socks5
         - http
         - local-forward
         - remote-forward
         - master
```
### SOCKS5
```
user@linux:~$ fullproxy socks5 --help
Usage of socks5:
  -address string
        Address to listen on. When "-slave" flag is set, is the IP of master to connect
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
Coming soon
### Forward
#### Local
```
user@linux:~$ fullproxy local-forward -help
Usage of local-forward:
  -forward-address string
        Address to forward the traffic received from master
  -forward-port string
        Port to forward the traffic received from master
  -master-address string
        Address of the master
  -master-port string
        Port of the master
```
#### Remote
```
user@linux:~$ fullproxy remote-forward -help
Usage of remote-forward:
  -local-address string
        Address accessible by master
  -localPort string
        Port of the address that is accessible by master
  -master-address string
        Address of the master
  -masterPort string
        Port of the master
```
### Master
```
user@linux:~$ fullproxy remote-forward -help
Usage of master:
  -address string
        Address to listen on. (default "0.0.0.0")
  -port string
        Port to listen on. (default "1080")
  -remote-address string
        Argument required to handle correctly the "remote-forward"
  -remote-port string
        Argument required to handle correctly the "remote-forward"
```
# Concepts
## Slave
Handles the proxying between a reverse connected (with encryption) proxy and the clients. In other words, it will receive the connections of the clients and will forward the traffic to the proxy that is reverse connected to it.
### How it works
1. It first binds to the address specified by the user.
2. Then accept the connection from the proxy server.
3. Finally, it proxy the traffic of all new incoming connections to the proxy server that was reverse connected to it in the second step.
In other words, is the proxy of another proxy but totally invisible for the client.
### Applications
This could be specially useful when you need to proxy a network that a machine have access to, but you can't bind with it
### Considerations
- The `master` protocol may loss some setup connections if it is extremely stressed, but it should work `just fine` if the connections where already made
# Installation
## Pre-compiled binaries
You can find pre-compiled binaries for windows and linux [Here](https://github.com/shoriwe/FullProxy/tree/master/build)
## Build from source code
Run this command:
```
go get github.com/shoriwe/FullProxy
```
# Suggestions
If you have any suggestion for new features, also leave them in the issue section or create the proper branch, add what do you want and request a pull request
