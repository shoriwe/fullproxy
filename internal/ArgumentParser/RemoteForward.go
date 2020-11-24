package ArgumentParser

import (
	"flag"
	"os"
	"time"
)

func RemotePortForwardArguments() (*string, *string, *string, *string, *int, *time.Duration, [2]string) {
	protocolFlagSet := flag.NewFlagSet("remote-forward", flag.ExitOnError)
	localHost := protocolFlagSet.String("local-host", "", "Host to bind by slave")
	localPort := protocolFlagSet.String("local-port", "", "Port to bind by slave")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	inboundWhitelist := protocolFlagSet.String("inbound-whitelist", "", "File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)")
	inboundBlacklist := protocolFlagSet.String("inbound-blacklist", "", "File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
	return localHost, localPort, masterHost, masterPort, tries, &timeout, [2]string{*inboundWhitelist, *inboundBlacklist}
}
