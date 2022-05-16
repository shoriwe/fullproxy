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
		listen string
		master string
	)
	httpCmd := flag.NewFlagSet("http", flag.ExitOnError)
	httpCmd.StringVar(&listen, "listen", "", "Address to listen for clients. Argument URL structure is 'network://host:port'")
	httpCmd.StringVar(&master, "master", "", "Listen address for master/slave communication. Argument URL structure is 'network://host:port'")
	parseError := httpCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	protocol := http2.NewHTTP(nil)
	log.Fatal(listeners.ServeHTTPHandler(listener, protocol, nil))
}
