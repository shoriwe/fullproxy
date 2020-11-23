package ArgumentParser

import (
	"flag"
	"os"
	"time"
)

func ParseMasterArguments() (*string, *string, *string, *string, *int, *time.Duration) {
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	host := protocolFlagSet.String("host", "0.0.0.0", "Host to listen on.")
	port := protocolFlagSet.String("port", "1080", "Port to listen on.")
	remoteHost := protocolFlagSet.String("forward-host", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	remotePort := protocolFlagSet.String("forward-port", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
	return host, port, remoteHost, remotePort, tries, &timeout
}
