# fullproxy

![build](https://img.shields.io/badge/build-passing-green)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/shoriwe/fullproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/fullproxy)](https://goreportcard.com/report/github.com/shoriwe/fullproxy)
[![Release](https://img.shields.io/github/release/shoriwe/fullproxy.svg?style=flat-square)](https://github.com/shoriwe/fullproxy/releases/latest)

**fullproxy** is a `bind` and `master/slave` based proxy toolkit.

![logo](logo/white_logo_color_background.jpg)

## Available protocols

- SOCKS5
- HTTP
- Port forward
- Translation from raw port to SOCKS5
- Raw port load balancer
- HTTP based load balancer

## Documentation

- [CLI](docs/cli.md)
- [YAML](docs/yaml.md)
- [Scripting](docs/scripting.md)
- [Pipes](docs/pipes.md)
- [Proxying](docs/proxy.md)
- [Filters](docs/filters.md)
- [Plasma programming language](https://shoriwe.github.io)

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
