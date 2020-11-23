package ArgumentParser

import (
	"flag"
	"os"
)

func RemotePortForwardArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("remote-forward", flag.ExitOnError)
	localHost := protocolFlagSet.String("local-host", "", "Host to bind by slave")
	localPort := protocolFlagSet.String("local-port", "", "Port to bind by slave")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return localHost, localPort, masterHost, masterPort
}
