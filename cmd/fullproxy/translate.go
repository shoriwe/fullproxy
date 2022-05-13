package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	pf_to_socks5 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/translation/pf-to-socks5"
	"log"
	"os"
)

func translate() {
	var (
		listen           string
		source           string
		sourceProtocol   string
		master           string
		target           string
		targetProtocol   string
		protocol         servers.Protocol
		listener         listeners.Listener
		newProtocolError error
		newListenerError error
	)
	translateCmd := flag.NewFlagSet("translate", flag.ExitOnError)
	translateCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	translateCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	translateCmd.StringVar(&source, "source", "", "Address of the source protocol")
	translateCmd.StringVar(&sourceProtocol, "sourceProtocol", "", "Protocol used by the source proxy")
	translateCmd.StringVar(&target, "target", "", "Target address.")
	translateCmd.StringVar(&targetProtocol, "targetProtocol", "", "Protocol used by the target proxy")
	parseError := translateCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	if listen == "" {
		printAndExit("no listen address provided", 1)
	}
	if source == "" {
		printAndExit("no source address provided", 1)
	}
	if sourceProtocol == "" {
		printAndExit("no source protocol provided", 1)
	}
	if target == "" {
		printAndExit("no target address provided", 1)
	}
	if targetProtocol == "" {
		printAndExit("no target protocol provided", 1)
	}
	if master != "" {
		listener, newListenerError = listeners.NewMaster("tcp", listen, nil, "tcp", master, nil)
	} else {
		listener, newListenerError = listeners.NewBindListener("tcp", listen, nil)
	}
	if newListenerError != nil {
		printAndExit(newListenerError.Error(), 1)
	}
	if sourceProtocol == "socks5" && targetProtocol == "forward" {
		protocol, newProtocolError = pf_to_socks5.NewForwardToSocks5("tcp", source, "", "", target)
	}
	if newProtocolError != nil {
		printAndExit(newProtocolError.Error(), 1)
	}
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
