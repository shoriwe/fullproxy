package main

import (
	"fmt"
	"github.com/shoriwe/FullProxy/internal/ArgumentParser"
	"github.com/shoriwe/FullProxy/internal/ProxiesSetup"
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
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
		host, port, username, password, slave, tries, timeout := ArgumentParser.ParseSocks5Arguments()
		ProxiesSetup.SetupSocks5(host, port, slave, username, password, *tries, *timeout)
	case "http":
		host, port, username, password, slave, tls := ArgumentParser.HTTPArguments()
		ProxiesSetup.SetupHTTP(host, port, slave, tls, username, password)
	case "local-forward":
		host, port, masterHost, masterPort, tries, timeout := ArgumentParser.LocalPortForwardArguments()
		ProxiesSetup.SetupLocalForward(host, port, masterHost, masterPort, *tries, *timeout)
	case "remote-forward":
		localHost, localPort, masterHost, masterPort, tries, timeout := ArgumentParser.RemotePortForwardArguments()
		ProxiesSetup.SetupRemoteForward(localHost, localPort, masterHost, masterPort, *tries, *timeout)
	case "master":
		masterHost, masterPort, remoteHost, remotePort, tries, timeout := ArgumentParser.ParseMasterArguments()
		if len(*remoteHost) > 0 && len(*remotePort) > 0 {
			ControllersSetup.MasterRemote(masterHost, masterPort, remoteHost, remotePort, *tries, *timeout)
		} else {
			ControllersSetup.MasterGeneral(masterHost, masterPort, *tries, *timeout)
		}
	case "translate":
		if numberOfArguments == 2 {
			_, _ = fmt.Fprintln(os.Stderr, "Try:\n", os.Args[0], " translate help")
			os.Exit(-1)
		}
		switch os.Args[2] {
		case "port_forward-socks5":
			bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort, tries, timeout := ArgumentParser.ParseForwardToSocks5Arguments()
			ProxiesSetup.SetupForwardSocks5(
				bindHost, bindPort,
				socks5Host, socks5Port,
				username, password,
				targetHost, targetPort,
				*tries, *timeout)
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
