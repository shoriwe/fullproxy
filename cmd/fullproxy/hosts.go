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
		listen string
		master string
	)
	hostsCmd := flag.NewFlagSet("hosts", flag.ExitOnError)
	hostsCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	hostsCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	parseError := hostsCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	protocol := http_hosts.NewHosts()
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
