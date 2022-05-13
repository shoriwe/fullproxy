package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	http_hosts "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http-hosts"
	"log"
	"os"
)

func hosts() {
	var (
		listen           string
		master           string
		listener         listeners.Listener
		newListenerError error
	)
	hostsCmd := flag.NewFlagSet("hosts", flag.ExitOnError)
	hostsCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	hostsCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	parseError := hostsCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	if listen == "" {
		printAndExit("no listen address provided", 1)
	}
	if master != "" {
		listener, newListenerError = listeners.NewMaster("tcp", listen, nil, "tcp", master, nil)
	} else {
		listener, newListenerError = listeners.NewBindListener("tcp", listen, nil)
	}
	if newListenerError != nil {
		printAndExit(newListenerError.Error(), 1)
	}
	protocol := http_hosts.NewHosts()
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
