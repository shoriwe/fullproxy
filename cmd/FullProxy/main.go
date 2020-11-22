package main

import (
	"fmt"
	"github.com/shoriwe/FullProxy/internal/ArgumentParser"
	"github.com/shoriwe/FullProxy/internal/ProxiesSetup"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers/Master"
	"os"
)

func main() {
	numberOfArguments := len(os.Args)
	if numberOfArguments == 1 {
		_, _ = fmt.Fprintln(os.Stderr, "Try:\n", os.Args[0], " help")
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "socks5":
		host, port, username, password, slave := ArgumentParser.ParseSocks5Arguments()
		ProxiesSetup.SetupSocks5(host, port, slave, &username, &password)
	case "http":
		host, port, username, password, slave, tls := ArgumentParser.HTTPArguments()
		ProxiesSetup.SetupHTTP(host, port, slave, tls, &username, &password)
	case "local-forward":
		host, port, masterHost, masterPort := ArgumentParser.LocalPortForwardArguments()
		ProxiesSetup.SetupLocalForward(host, port, masterHost, masterPort)
	case "remote-forward":
		localHost, localPort, masterHost, masterPort := ArgumentParser.RemotePortForwardArguments()
		ProxiesSetup.SetupRemoteForward(localHost, localPort, masterHost, masterPort)
	case "master":
		masterHost, masterPort, remoteHost, remotePort := ArgumentParser.ParseMasterArguments()
		Master.Master(masterHost, masterPort, remoteHost, remotePort)
	case "translate":
		if numberOfArguments == 2 {
			_, _ = fmt.Fprintln(os.Stderr, "Try:\n", os.Args[0], " translate help")
			os.Exit(-1)
		}
		switch os.Args[2] {
		case "port_forward-socks5":
			bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort := ArgumentParser.ParseForwardToSocks5Arguments()
			ProxiesSetup.SetupForwardSocks5(bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort)
		case "help":
			ArgumentParser.ShowTranslateHelpMessage()
		default:
			_, _ = fmt.Fprintln(os.Stderr, "Unknown command\nTry: ", os.Args[0], " translate help")
			os.Exit(-1)
		}
	case "help":
		ArgumentParser.ShowGeneralHelpMessage()
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Unknown command\nTry: ", os.Args[0], " help")
		os.Exit(-1)
	}
}
