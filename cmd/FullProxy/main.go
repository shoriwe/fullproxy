package main

import (
	"fmt"
	"github.com/shoriwe/FullProxy/internal/ArgumentParser"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/internal/ProxiesSetup"
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
		host, port, username, password, slave, tries, timeout, inboundLists, outboundLists := ArgumentParser.ParseSocks5Arguments()
		ProxiesSetup.SetupSocks5(host, port, slave, username, password, tries, timeout, inboundLists, outboundLists)
	case "http":
		host, port, username, password, slave, tls, inboundLists, outboundLists := ArgumentParser.HTTPArguments()
		ProxiesSetup.SetupHTTP(host, port, slave, tls, username, password, inboundLists, outboundLists)
	case "local-forward":
		host, port, masterHost, masterPort, tries, timeout, inboundLists := ArgumentParser.LocalPortForwardArguments()
		ProxiesSetup.SetupLocalForward(host, port, masterHost, masterPort, tries, timeout, inboundLists)
	case "remote-forward":
		localHost, localPort, masterHost, masterPort, tries, timeout, inboundLists := ArgumentParser.RemotePortForwardArguments()
		ProxiesSetup.SetupRemoteForward(localHost, localPort, masterHost, masterPort, tries, timeout, inboundLists)
	case "master":
		masterHost, masterPort, remoteHost, remotePort, tries, timeout, inboundLists := ArgumentParser.ParseMasterArguments()
		if len(*remoteHost) > 0 && len(*remotePort) > 0 {
			PipesSetup.RemoteForwardMaster(masterHost, masterPort, remoteHost, remotePort, tries, timeout)
		} else {
			PipesSetup.GeneralMaster(masterHost, masterPort, tries, timeout, inboundLists)
		}
	case "translate":
		if numberOfArguments == 2 {
			_, _ = fmt.Fprintln(os.Stderr, "Try:\n", os.Args[0], " translate help")
			os.Exit(-1)
		}
		switch os.Args[2] {
		case "port_forward-socks5":
			bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort, tries, timeout, inboundLists := ArgumentParser.ParseForwardToSocks5Arguments()
			ProxiesSetup.SetupForwardSocks5(
				bindHost, bindPort,
				socks5Host, socks5Port,
				username, password,
				targetHost, targetPort,
				tries, timeout,
				inboundLists)
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
