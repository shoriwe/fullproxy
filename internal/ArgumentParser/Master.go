package ArgumentParser

import (
	"flag"
	"os"
)

func ParseMasterArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	host := protocolFlagSet.String("host", "0.0.0.0", "Host to listen on.")
	port := protocolFlagSet.String("port", "1080", "Port to listen on.")
	remoteHost := protocolFlagSet.String("forward-host", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	remotePort := protocolFlagSet.String("forward-port", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return host, port, remoteHost, remotePort
}
