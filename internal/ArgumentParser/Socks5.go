package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func ParseSocks5Arguments() (*string, *string, []byte, []byte, *bool, *int, *time.Duration, [2]string, [2]string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	host := protocolFlagSet.String("host", "", "Host to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "1080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 1080")
	username := protocolFlagSet.String("username-auth", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string, this auth method will ignore any other supplied")
	password := protocolFlagSet.String("password-auth", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string, this auth method will ignore any other supplied")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	inboundWhitelist := protocolFlagSet.String("inbound-whitelist", "", "File with a host per line. Allowed incoming connections to the proxy (ignored in slave  mode and when inbound-blacklist is set)")
	inboundBlacklist := protocolFlagSet.String("inbound-blacklist", "", "File with a host per line. Denied incoming connections to the proxy (ignored in slave mode and when inbound-whitelist is set)")
	outboundWhitelist := protocolFlagSet.String("outbound-whitelist", "", "File with a host per line. Allowed outgoing connections (ignored when outbound-blacklist is set)")
	outboundBlacklist := protocolFlagSet.String("outbound-blacklist", "", "File with a host per line. Denied outgoing connections (ignored when outbound-whitelist is set)")
	commandAuth := protocolFlagSet.String("command-auth", "", "Command with it's default args to pass the Username and Password received from clients, please notice that ExitCode = 0 will mean that the login was successful, any other way it not and the username and password will be passed as base64 encoded arguments to it, this auth method will ignore any other supplied")
	databaseAuth := protocolFlagSet.String("database-auth", "", "Path to the SQLite3 database generated with the 'database create' command and filled with the 'database user add' command, this auth method will ignore any other supplied")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
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
	return host, port, []byte(*username), []byte(*password), slave, tries, &timeout, [2]string{*inboundWhitelist, *inboundBlacklist}, [2]string{*outboundWhitelist, *outboundBlacklist}, commandAuth, databaseAuth
}
