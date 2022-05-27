package main

import (
	"crypto/tls"
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	port_forward "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/port-forward"
	"log"
	"net/url"
	"os"
)

func forward() {
	var (
		listen  string
		master  string
		target  string
		dialTLS bool
	)
	forwardCmd := flag.NewFlagSet("forward", flag.ExitOnError)
	forwardCmd.StringVar(&listen, "listen", "", "Address to listen for clients. Argument URL structure is 'network://host:port'")
	forwardCmd.StringVar(&master, "master", "", "Listen address for master/slave communication. Argument URL structure is 'network://host:port'")
	forwardCmd.StringVar(&target, "target", "", "Target forward address. Argument URL structure is 'network://host:port'")
	forwardCmd.BoolVar(&dialTLS, "dial-tls", false, "Dial connection will use TLS")
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
	var tlsConfig *tls.Config
	if dialTLS {
		tlsConfig = &tls.Config{InsecureSkipVerify: dialTLS}
	}
	protocol := port_forward.NewForward(targetUrl.Scheme, targetUrl.Host, tlsConfig)
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
