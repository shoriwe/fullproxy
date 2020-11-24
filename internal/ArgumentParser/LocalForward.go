package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func LocalPortForwardArguments() (*string, *string, *string, *string, *int, *time.Duration, [2]string) {
	protocolFlagSet := flag.NewFlagSet("local-forward", flag.ExitOnError)
	forwardHost := protocolFlagSet.String("forward-host", "", "Host to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forward-port", "", "Port to forward the traffic received from master")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	inboundWhitelist := protocolFlagSet.String("inbound-whitelist", "", "File with a host per line. Allowed incoming connections to the proxy (ignored when inbound-blacklist is set)")
	inboundBlacklist := protocolFlagSet.String("inbound-blacklist", "", "File with a host per line. Denied incoming connections to the proxy (ignored when inbound-whitelist is set)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
	if len(*inboundWhitelist) > 0 && len(*inboundBlacklist) > 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot use Inbound Whitelist with an Inbound Blacklist")
	}
	return forwardHost, forwardPort, masterHost, masterPort, tries, &timeout, [2]string{*inboundWhitelist, *inboundBlacklist}
}
