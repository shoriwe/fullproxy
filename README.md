# fullproxy

![build](https://img.shields.io/badge/build-passing-green)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/shoriwe/fullproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/fullproxy/v3)](https://goreportcard.com/report/github.com/shoriwe/fullproxy/v3)
[![Release](https://img.shields.io/github/release/shoriwe/fullproxy.svg?style=flat-square)](https://github.com/shoriwe/fullproxy/releases/latest)

**fullproxy** is a `listen port` and `master/slave` based proxy toolkit.

![logo](logo/white_logo_color_background.jpg)

## Available proxy protocols

- SOCKS5
- HTTP
- Port forward
- Translation from raw port to SOCKS5.
- Raw port reverse proxy and load balancer
- HTTP reverse proxy and load balancer

## Quick preview

### Listen port communication

This is the classic bind to port and handle connections way. Useful when the user have the necessary permissions in the machine that has access to the targeted networks.

### Master/Slave communication

This protocol is useful to forward traffic from the networks a machine have access to but cannot listen to, only permitting outgoing specifically to our master server.

### CLI

#### Help

```shell
fullproxy help
```

```
Usage: fullproxy COMMAND [ARGUMENTS]

Available commands:

- help:                 prints this help.
- slave:                Connects to master server.
- socks5:               Starts a SOCKS5 server.
- http:                 Starts a HTTP proxy server.
- forward:              Starts a port forward proxy server.
- translate:    Translate a proxy protocol to another to proxy protocol.
- reverse:              Starts a raw reverse proxy.
- config:               Start serving the server configured in the targeted yaml file.
```

#### Slave

Slave is protocol independent, it can be used with any proxy protocol that specifies the `master MASTER_IP` flag.

##### Help message

```shell
fullproxy slave -h
```

```
Usage of translate:
  -master string
        Address of master server. Argument URL structure is 'network://host:port'
```

##### Usage example

- Master machine

```shell
fullproxy socks5 -listen tcp://0.0.0.0:9050 -master tcp://192.168.1.33:9051
```

- Slave machine

```shell
fullproxy slave -master tcp://192.168.1.33:9051
```

#### SOCKS5

##### Help message

```shell
fullproxy socks5 -h
```

```
Usage of socks5:
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

##### Usage example

```shell
fullproxy socks5 -listen tcp://192.168.1.33:9050
```

#### HTTP

##### Help message

```shell
fullproxy http -h
```

```
Usage of http:
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

##### Usage example

```shell
fullproxy http -listen tcp://192.168.1.33:9050
```

#### Port forwarding

Receive connections and forward the traffic to a predefined target.

##### Help message

```shell
fullproxy forward -h
```

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

##### Usage example

###### Basic usage

```shell
fullproxy forward -listen tcp://127.0.0.1:80 -target tcp://192.168.1.34:80
```

###### Receive raw connection, forward to TLS port

```shell
fullproxy forward -listen tcp://127.0.0.1:80 -dial-tls -target tcp://google.com:443
```

#### Protocol translation

Translate one proxy protocol to another.

##### Help message

```shell
fullproxy translate -h
```

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

##### Usage example

###### Raw port to SOCKS5 target

```shell
fullproxy translate -listen tcp://127.0.0.1:443 -source tcp://192.168.1.33:9050/socks5 -target tcp://google.com:443/forward
```

#### Raw reverse proxy and load balancer

Receive connection and redirect traffic to a any host specified in a pool.

##### Help message

```shell
fullproxy reverse -h
```

```
Usage of reverse:
  -add value
        Add target to the pool load balancer
  -listen string
        Address to listen for clients. Argument URL structure is 'network://host:port'
  -master string
        Listen address for master/slave communication. Argument URL structure is 'network://host:port'
```

##### Usage example

```shell
fullproxy reverse -listen tcp://0.0.0.0:80 -add tcp://192.168.1.33:80 -add tcp://192.168.1.34:80 -add tcp://192.168.1.35:80
```

### Advance usage

`fullproxy` have even more features that can be configured using `config` subcommand which loads and prewritten config YAML file.

#### Help message

```shell
fullproxy config
```

 ```
 Usage: fullproxy config YAML_CONFIG
 ```

It also prints the [docs/yaml.md](docs/yaml.md) file with a sample configuration file.

```yaml
init-order:
  - LISTENER_NAME
  - ...
drivers: # This field is optional.
  DRIVER_NAME: /PATH/TO/PLASMA/SCRIPT
services:
  SERVICE_NAME:
    log: /PATH/TO/FILE/TO/DATA
    sniff:
      incoming: /PATH/TO/FILE/WITH/INCOMING/TRAFFIC
      outgoing: /PATH/TO/FILE/WITH/OUTGOING/TRAFFIC
    listener:
      # Mandatory by all types of listeners
      type: basic | master | slave
      network: tcp | unix
      address: HOST:PORT | /PATH/TO/UNIX/SOCK
      tls: # Ignore to generate a self signed cert
        - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY
        - ...
      # Available for all types of protocol
      filters:
        inbound: DRIVER_NAME # Ignore to no filer
        outbound: DRIVER_NAME # Ignore to no filer
        listen: DRIVER_NAME # Ignore to no filer
        accept: DRIVER_NAME # Ignore to no filer

      # Mandatory by master and slave
      master-network: tcp | unix
      master-address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory by master
      master-tls: # Ignore to generate a self signed cert.
        - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY

      # Mandatory by slave
      slave-trust: true | false

    protocol: # Used only when type is basic | master
      # Mandatory
      type: socks5|http|reverse-raw|reverse-http|forward|translate

      # Only for socks5 and http
      authentication: DRIVE_NAME # Ignore to no auth

      # Mandatory by forward
      dial-tls:
        trust: true|false
        certificate: /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY

      # Mandatory by forward and translate
      target-network: tcp | unix
      target-address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory by translate
      proxy-network: tcp | unix
      proxy-address: HOST:PORT | /PATH/TO/UNIX/SOCK
      translation: socks5:forward # Currently only supported
      credentials: USERNAME:PASSWORD

      # Mandatory for reverse-raw
      raw-hosts:
        NAME:
          tls: # Do not set to use raw sockets
            trust: true | false
            certificates:
              - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY
              - ...
          network: tcp | unix
          address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory for reverse-http
      http-hosts:
        HOSTNAME:
          uri: /wanted/uri
          response-headers: # Headers to in inject in the response to the client
            KEY: VALUE
          request-headers: # Headers to in inject in the request to the server
            KEY: VALUE
          pool: # Load balancing pool
            NAME:
              websocket-read-buffer-size: NUMERIC # Ignore both buffer size settings to not forward websocket traffic
              websocket-write-buffer-size: NUMERIC # Ignore both buffer size settings to not forward websocket traffic
              tls: # Do not set to use raw sockets
                trust: true | false
                certificates:
                  - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY
                  - ...
              scheme: http | https | ws | wss # This is used only for websocket connections 
              uri: URI
              network: tcp | unix
              address: HOST:PORT | /PATH/TO/UNIX/SOCK
```

#### Example usage

Check the [sample-config/configs/sample.yaml](sample-config/configs/sample.yaml) config file for reference.

```shell
fullproxy config sample-config/configs/sample.yaml
```

```
2022/05/29 15:47:37 Loading drivers
2022/05/29 15:47:37 Drivers loaded
2022/05/29 15:47:37 Starting listeners
2022/05/29 15:47:37 - http-master
2022/05/29 15:47:37 - socks5
2022/05/29 15:47:37 - google
2022/05/29 15:47:37 - http-slave
2022/05/29 15:47:37 - forward
2022/05/29 15:47:37 - reverse-http
2022/05/29 15:47:37 - reverse-raw
2022/05/29 15:47:37 Listeners started
2022/05/29 15:47:37 forward Started
2022/05/29 15:47:37 socks5 Started
2022/05/29 15:47:37 reverse-http Started
2022/05/29 15:47:37 google Started
2022/05/29 15:47:37 reverse-raw Started
2022/05/29 15:47:37 http-master Started
...
```

## Detailed documentation

- [CLI](docs/cli.md)
- [YAML](docs/yaml.md)
- [Scripting](docs/scripting.md)
- [Pipes](docs/pipes.md)
- [Proxying](docs/proxy.md)
- [Filters](docs/filters.md)
- [Plasma programming language](https://shoriwe.github.io/plasma/index.html)
- [pkg.go.dev](https://pkg.go.dev/github.com/shoriwe/fullproxy/v3)

## Installation

### Pre-compiled binaries

You can find pre-compiled binaries for windows and Linux [Here](https://github.com/shoriwe/fullproxy/releases)

### Build from source code

#### Go tool based

```shell
go install github.com/shoriwe/fullproxy/v3@latest
```

#### Git clone based

```shell
git clone https://github.com/shoriwe/fullproxy
cd fullproxy
go build -mod vendor -o fullproxy
```
