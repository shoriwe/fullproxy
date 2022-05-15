package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	pf_to_socks5 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/translation/pf-to-socks5"
	"log"
	"net/url"
	"os"
)

func selectTranslateProtocol(source, target string) (servers.Protocol, error) {
	if source == "" {
		return nil, errors.New("no source address provided")
	}
	sourceUrl, sourceUrlError := url.Parse(source)
	if sourceUrlError != nil {
		return nil, sourceUrlError
	}
	if target == "" {
		return nil, errors.New("no target address provided")
	}
	targetUrl, targetUrlError := url.Parse(target)
	if targetUrlError != nil {
		return nil, targetUrlError
	}
	if sourceUrl.Path == "/socks5" && targetUrl.Path == "/forward" {
		return pf_to_socks5.NewForwardToSocks5(
			sourceUrl.Scheme, sourceUrl.Host,
			sourceUrl.User,
			targetUrl.Scheme, targetUrl.Host,
		)
	} else {
		return nil, fmt.Errorf("unknown translation between %s and %s", sourceUrl.Path, targetUrl.Path)
	}
}

func translate() {
	var (
		listen string
		source string
		master string
		target string
	)
	translateCmd := flag.NewFlagSet("translate", flag.ExitOnError)
	translateCmd.StringVar(&listen, "listen", "", "Address to listen for clients. Argument URL structure is 'network://host:port'")
	translateCmd.StringVar(&master, "master", "", "Listen address for master/slave communication. Argument URL structure is 'network://host:port'")
	translateCmd.StringVar(&source, "source", "", "Address of the source proxy. Argument URL structure is 'network://host:port/protocol'")
	translateCmd.StringVar(&target, "target", "", "Address of target proxy. Argument URL structure is 'network://host:port/protocol'")
	parseError := translateCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	protocol, protocolError := selectTranslateProtocol(
		source,
		target,
	)
	if protocolError != nil {
		printAndExit(protocolError.Error(), 1)
	}
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
