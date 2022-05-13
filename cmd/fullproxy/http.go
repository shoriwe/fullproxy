package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	http2 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/http"
	"log"
	"os"
)

func http() {
	var (
		listen           string
		master           string
		listener         listeners.Listener
		newListenerError error
	)
	httpCmd := flag.NewFlagSet("http", flag.ExitOnError)
	httpCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	httpCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	parseError := httpCmd.Parse(os.Args[2:])
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
	protocol := http2.NewHTTP(nil)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
