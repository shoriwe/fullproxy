package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
)

func HTTPArguments() (*string, *string, []byte, []byte, *bool, *bool, [2]string, [2]string) {
	protocolFlagSet := flag.NewFlagSet("http", flag.ExitOnError)
	host := protocolFlagSet.String("host", "", "Host to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "8080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 8080")
	username := protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	tls := protocolFlagSet.Bool("tls", false, "Use HTTPS")
	inboundWhitelist := protocolFlagSet.String("inbound-whitelist", "", "File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)")
	inboundBlacklist := protocolFlagSet.String("inbound-blacklist", "", "File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)")
	outboundWhitelist := protocolFlagSet.String("outbound-whitelist", "", "File with a host per line. Allowed outgoing connections (ignored when outbound-blacklist is set)")
	outboundBlacklist := protocolFlagSet.String("outbound-blacklist", "", "File with a host per line. Denied outgoing connections (ignored when outbound-whitelist is set)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	if len(*host) == 0 {
		if *slave {
			*host = "127.0.0.1"
		} else {
			*host = "0.0.0.0"
		}
	}
	if len(*inboundWhitelist) > 0 && len(*inboundBlacklist) > 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot use Inbound Whitelist with an Inbound Blacklist")
	}
	if len(*outboundWhitelist) > 0 && len(*outboundBlacklist) > 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Cannot use Outbound Whitelist with an Outbound Blacklist")
	}
	return host, port, []byte(*username), []byte(*password), slave, tls, [2]string{*inboundWhitelist, *inboundBlacklist}, [2]string{*outboundWhitelist, *outboundBlacklist}
}
