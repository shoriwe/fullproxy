package main

import (
	"github.com/shoriwe/FullProxy/internal/ArgumentParser"
	"github.com/shoriwe/FullProxy/internal/Database"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/internal/ProxiesSetup"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"os"
)

func main() {
	numberOfArguments := len(os.Args)
	Templates.Exit("Try:\n\t%s help\n", 2, numberOfArguments, os.Args[0])
	switch os.Args[1] {
	case "socks5":
		host, port, username, password, slave, tries, timeout, inboundLists, outboundLists, commandAuth, databaseAuth := ArgumentParser.ParseSocks5Arguments()
		ProxiesSetup.SetupSocks5(host, port, slave, username, password, tries, timeout, inboundLists, outboundLists, commandAuth, databaseAuth)
	case "http":
		host, port, username, password, slave, tls, inboundLists, outboundLists, commandAuth, databaseAuth := ArgumentParser.HTTPArguments()
		ProxiesSetup.SetupHTTP(host, port, slave, tls, username, password, inboundLists, outboundLists, commandAuth, databaseAuth)
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
		Templates.Exit("Try:\n\t%s translate help\n", 3, numberOfArguments, os.Args[0])
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
			Templates.Exit("Unknown command\nTry:\n\t%s translate help\n", 1, 0, os.Args[0])
		}
	case "database":
		Templates.Exit("No command supplied\nTry:\n\t%s database help\n", 3, numberOfArguments, os.Args[0])
		switch os.Args[2] {
		case "create":
			Templates.Exit("Usage:\n\t%s database create DATABASE_FILE\n", 4, numberOfArguments, os.Args[0])
			Database.Create(os.Args[3])
		case "user":
			Templates.Exit("No command received\nTry:\n\t%s database user help\n", 4, numberOfArguments, os.Args[0])
			switch os.Args[3] {
			case "add":
				Templates.Exit("Usage:\n\t%s database user add DATABASE_FILE USERNAME PASSWORD", 7, numberOfArguments, os.Args[0])
				Database.Add(os.Args[4], os.Args[5], os.Args[6])
			case "delete":
				Templates.Exit("Usage:\n\t%s database user delete DATABASE_FILE USERNAME", 6, numberOfArguments, os.Args[0])
				Database.Delete(os.Args[4], os.Args[5])
			case "update":
				Templates.Exit("Usage:\n\t%s database user update DATABASE_FILE USERNAME NEW_PASSWORD", 7, numberOfArguments, os.Args[0])
				Database.Update(os.Args[4], os.Args[5], os.Args[6])
			case "help":
				ArgumentParser.ShowDatabaseUserHelpMessage()
			default:
				Templates.Exit("Unknown command\nTry: %s database user help", 1, 0, os.Args[0])
			}
		case "help":
			ArgumentParser.ShowDatabaseHelpMessage()
		default:
			Templates.Exit("Unknown command\nTry: %s database help", 1, 0, os.Args[0])
		}
	case "help":
		ArgumentParser.ShowGeneralHelpMessage()
	default:
		Templates.Exit("Unknown command\nTry: %s help", 1, 0, os.Args[0])
	}
}
