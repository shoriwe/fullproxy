package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	port_forward "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	"log"
	"os"
)

func forward() {
	var (
		listen           string
		master           string
		target           string
		listener         listeners.Listener
		newListenerError error
	)
	forwardCmd := flag.NewFlagSet("forward", flag.ExitOnError)
	forwardCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	forwardCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	forwardCmd.StringVar(&target, "target", "", "Target address to redirect the traffic.")
	parseError := forwardCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	if listen == "" {
		printAndExit("no listen address provided", 1)
	}
	if target == "" {
		printAndExit("no target address provided", 1)
	}
	if master != "" {
		listener, newListenerError = listeners.NewMaster("tcp", listen, nil, "tcp", master, nil)
	} else {
		listener, newListenerError = listeners.NewBindListener("tcp", listen, nil)
	}
	if newListenerError != nil {
		printAndExit(newListenerError.Error(), 1)
	}
	protocol := port_forward.NewForward(target)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
