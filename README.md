# fullproxy

[![Release](https://github.com/shoriwe/fullproxy/actions/workflows/release.yml/badge.svg)](https://github.com/shoriwe/fullproxy/actions/workflows/release.yml)
[![Test](https://github.com/shoriwe/fullproxy/actions/workflows/test.yml/badge.svg)](https://github.com/shoriwe/fullproxy/actions/workflows/test.yml)
[![Versioning](https://github.com/shoriwe/fullproxy/actions/workflows/version.yml/badge.svg)](https://github.com/shoriwe/fullproxy/actions/workflows/version.yml)
[![codecov](https://codecov.io/gh/shoriwe/fullproxy/branch/master/graph/badge.svg?token=WQSZVR7YT7)](https://codecov.io/gh/shoriwe/fullproxy)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shoriwe/fullproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoriwe/fullproxy/v4)](https://goreportcard.com/report/github.com/shoriwe/fullproxy/v4)

<img src="logo/white_logo_color_background.jpg" alt="logo" style="zoom: 25%;" />

## Installing

### Pre-compiled binaries 

You can find pre-compiled binaries in [release](releases/latest)

### Go install

Can compile from source with:

```shell
go install github.com/shoriwe/fullproxy/v4@latest
```

### Cloning repository

- Clone repository

```shell
git clone https://github.com/shoriwe/fullproxy
```

- Build

```shell
cd fullproxy && go build .
```

## Preview

### Compose

See [Compose](docs/Compose.md) for more information about **compose contracts**.

```shell
fullproxy compose ./fullproxy-compose.yaml
```

## Documentation

| File                                                     | Description                                       |
| -------------------------------------------------------- | ------------------------------------------------- |
| [Circuits](docs/Circuits.md)                             | Documentation about how circuits work             |
| [CLI](docs/CLI.md)                                       | Documentation of the CLI tool                     |
| [Compose](docs/Compose.md)                               | Documentation about the **compose** specification |
| [Continuous integration](docs/Continuous%20integration.md) | Documentation of the CI                           |

## Coverage

| [![codecov](https://codecov.io/gh/shoriwe/fullproxy/branch/master/graphs/sunburst.svg?token=WQSZVR7YT7)](https://github.com/shoriwe/fullproxy) | [![codecov](https://codecov.io/gh/shoriwe/fullproxy/branch/master/graphs/tree.svg?token=WQSZVR7YT7)](https://github.com/shoriwe/fullproxy) |
| :----------------------------------------------------------: | :----------------------------------------------------------: |

