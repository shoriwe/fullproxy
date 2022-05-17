package main

import (
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	port_forward "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	"log"
	"net/url"
	"os"
)

func forward() {
	var (
		listen string
		master string
		target string
	)
	forwardCmd := flag.NewFlagSet("forward", flag.ExitOnError)
	forwardCmd.StringVar(&listen, "listen", "", "Address to listen for clients. Argument URL structure is 'network://host:port'")
	forwardCmd.StringVar(&master, "master", "", "Listen address for master/slave communication. Argument URL structure is 'network://host:port'")
	forwardCmd.StringVar(&target, "target", "", "Target forward address. Argument URL structure is 'network://host:port'")
	parseError := forwardCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	if target == "" {
		printAndExit("no target address provided", 1)
	}
	targetUrl, parseTargetUrlError := url.Parse(target)
	if parseTargetUrlError != nil {
		printAndExit(parseTargetUrlError.Error(), 1)
	}
	protocol := port_forward.NewForward(targetUrl.Scheme, targetUrl.Host)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
