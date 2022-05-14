package main

import (
	"flag"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	pf_to_socks5 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/translation/pf-to-socks5"
	"log"
	"net/url"
	"os"
)

func translate() {
	var (
		listen           string
		source           string
		sourceProtocol   string
		master           string
		target           string
		targetProtocol   string
		protocol         servers.Protocol
		newProtocolError error
	)
	translateCmd := flag.NewFlagSet("translate", flag.ExitOnError)
	translateCmd.StringVar(&listen, "listen", "", "Address to listen for clients")
	translateCmd.StringVar(&master, "master", "", "Listen address for master/slave communication.")
	translateCmd.StringVar(&source, "source", "", "Address of the source protocol")
	translateCmd.StringVar(&sourceProtocol, "sourceProtocol", "", "Protocol used by the source proxy")
	translateCmd.StringVar(&target, "target", "", "Target address.")
	translateCmd.StringVar(&targetProtocol, "targetProtocol", "", "Protocol used by the target proxy")
	parseError := translateCmd.Parse(os.Args[2:])
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listener, listenError := createListener(listen, master)
	if listenError != nil {
		printAndExit(listenError.Error(), 1)
	}
	if source == "" {
		printAndExit("no source address provided", 1)
	}
	sourceUrl, sourceUrlError := url.Parse(source)
	if sourceUrlError != nil {
		printAndExit(sourceUrlError.Error(), 1)
	}
	if target == "" {
		printAndExit("no target address provided", 1)
	}
	targetUrl, targetUrlError := url.Parse(target)
	if targetUrlError != nil {
		printAndExit(targetUrlError.Error(), 1)
	}
	if sourceProtocol == "socks5" && targetProtocol == "forward" {
		protocol, newProtocolError = pf_to_socks5.NewForwardToSocks5(
			sourceUrl.Scheme, sourceUrl.Host,
			"", "",
			targetUrl.Scheme, targetUrl.Host,
		)
	} else {
		printAndExit(fmt.Sprintf("unknown translation between %s and %s", sourceProtocol, targetProtocol), 1)
	}
	if newProtocolError != nil {
		printAndExit(newProtocolError.Error(), 1)
	}
	log.Fatal(listeners.Serve(listener, protocol, nil))
}
