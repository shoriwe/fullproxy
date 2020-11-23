package ArgumentParser

import (
	"flag"
	"os"
	"time"
)

func LocalPortForwardArguments() (*string, *string, *string, *string, *int, *time.Duration) {
	protocolFlagSet := flag.NewFlagSet("local-forward", flag.ExitOnError)
	forwardHost := protocolFlagSet.String("forward-host", "", "Host to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forward-port", "", "Port to forward the traffic received from master")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	timeout := time.Duration(*rawTimeout) * time.Second
	return forwardHost, forwardPort, masterHost, masterPort, tries, &timeout
}