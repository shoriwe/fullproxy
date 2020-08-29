package ArgumentParser


import (
	"flag"
	"fmt"
	"os"
)


func PortForwardArguments() (*string, *string, *string, *string){
	protocolFlagSet := flag.NewFlagSet("forward", flag.ExitOnError)
	forwardAddress := protocolFlagSet.String( "forwardAddress", "", "Address to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forwardPort", "", "Port to forward the traffic received from master")
	masterAddress := protocolFlagSet.String("masterAddress", "", "Address of the master")
	masterPort := protocolFlagSet.String("masterPort", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return forwardAddress, forwardPort, masterAddress, masterPort
}


func ShowGeneralHelpMessage(){
	_, _ = fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " PROTOCOL *FLAGS\nProtocols available:\n\t - socks5\n\t - http\n\t - forward\n\t - master")
}


func ParseSocks5Arguments() (*string, *string, []byte, []byte, *bool) {
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	address := protocolFlagSet.String( "address", "", "Address to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "1080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 1080")
	username := protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	_ = protocolFlagSet.Parse(os.Args[2:])
	if len(*address) == 0{
		if *slave{
			*address = "127.0.0.1"
		} else {
			*address = "0.0.0.0"
		}
	}
	return address, port, []byte(*username), []byte(*password), slave
}


func ParseMasterArguments() (*string, *string){
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	address := protocolFlagSet.String( "address", "0.0.0.0", "Address to listen on.")
	port := protocolFlagSet.String( "port", "1080", "Port to listen on.")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return address, port
}

