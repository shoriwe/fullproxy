package ArgumentParser

import (
	"flag"
	"os"
)

func LocalPortForwardArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("local-forward", flag.ExitOnError)
	forwardHost := protocolFlagSet.String("forward-host", "", "Host to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forward-port", "", "Port to forward the traffic received from master")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return forwardHost, forwardPort, masterHost, masterPort
}