package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	"log"
	"os"
)

func socks5() {
	var (
		listen           string
		master           string
		listener         listeners.Listener
		newListenerError error
	)
	socks5Cmd := flag.NewFlagSet("socks5", flag.ExitOnError)
	socks5Cmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	socks5Cmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	parseError := socks5Cmd.Parse(os.Args[2:])
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
	protocol := socks52.NewSocks5(nil)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
