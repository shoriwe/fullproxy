package main

import (
	"fmt"
	"os"
)

const (
	helpMessage = `Usage: fullproxy COMMAND [ARGUMENTS]

Available commands:

- help:			prints this help.
- slave:		Connects to master server.
- socks5:		Starts a SOCKS5 server.
- http:			Starts a HTTP proxy server.
- forward:		Starts a port forward proxy server.
- translate:	Translate a proxy protocol to another to proxy protocol.
- reverse:		Starts a raw reverse proxy.
- config:		Start serving the server configured in the targeted yaml file.`
)

func printAndExit(msg string, code int) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func main() {
	if len(os.Args) == 1 {
		printAndExit(helpMessage, 0)
	}
	switch os.Args[1] {
	case "help":
		printAndExit(helpMessage, 0)
	case "slave":
		slave()
	case "socks5":
		socks5()
	case "http":
		http()
	case "forward":
		forward()
	case "translate":
		translate()
	case "reverse":
		reverse()
	case "config":
		config()
	default:
		printAndExit(fmt.Sprintf("Unknown command '%s'", os.Args[1]), 1)
	}
}
