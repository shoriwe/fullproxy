package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func ParseMasterArguments() (*string, *string, *string, *string, *int, *time.Duration, [2]string) {
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	host := protocolFlagSet.String("host", "0.0.0.0", "Host to listen on.")
	port := protocolFlagSet.String("port", "1080", "Port to listen on.")
	remoteHost := protocolFlagSet.String("forward-host", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	remotePort := protocolFlagSet.String("forward-port", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	inboundWhitelist := protocolFlagSet.String("inbound-whitelist", "", "File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)")
	inboundBlacklist := protocolFlagSet.String("inbound-blacklist", "", "File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
	if len(*inboundWhitelist) > 0 && len(*inboundBlacklist) > 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot use Inbound Whitelist with an Inbound Blacklist")
	}
	return host, port, remoteHost, remotePort, tries, &timeout, [2]string{*inboundWhitelist, *inboundBlacklist}
}
