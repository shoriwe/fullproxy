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
		listen string
		master string
	)
	socks5Cmd := flag.NewFlagSet("socks5", flag.ExitOnError)
	socks5Cmd.StringVar(&listen, "listen", "", "Address to listen for clients. Argument URL structure is 'network://host:port'")
	socks5Cmd.StringVar(&master, "master", "", "Listen address for master/slave communication. Argument URL structure is 'network://host:port'")
	parseError := socks5Cmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	protocol := socks52.NewSocks5(nil)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
