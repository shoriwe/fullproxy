package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	reverse2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/reverse"
	"log"
	"os"
)

func reverse() {
	var (
		pool             stringSlice
		listen           string
		master           string
		listener         listeners.Listener
		newListenerError error
	)
	reverseCmd := flag.NewFlagSet("reverse", flag.ExitOnError)
	reverseCmd.Var(&pool, "pool", "List of targets used by the load balancer.")
	reverseCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	reverseCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	parseError := reverseCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	if listen == "" {
		printAndExit("no listen address provided", 1)
	}
	if len(pool.contents) == 0 {
		printAndExit("no pool targets provided", 1)
	}
	if master != "" {
		listener, newListenerError = listeners.NewMaster("tcp", listen, nil, "tcp", master, nil)
	} else {
		listener, newListenerError = listeners.NewBindListener("tcp", listen, nil)
	}
	if newListenerError != nil {
		printAndExit(newListenerError.Error(), 1)
	}
	protocol := reverse2.NewRaw(pool.contents)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
