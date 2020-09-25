package ArgumentParser


import (
	"flag"
	"fmt"
	"os"
)


func LocalPortForwardArguments() (*string, *string, *string, *string){
	protocolFlagSet := flag.NewFlagSet("local-forward", flag.ExitOnError)
	forwardAddress := protocolFlagSet.String( "forward-address", "", "Address to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forward-port", "", "Port to forward the traffic received from master")
	masterAddress := protocolFlagSet.String("master-address", "", "Address of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return forwardAddress, forwardPort, masterAddress, masterPort
}


func RemotePortForwardArguments() (*string, *string, *string, *string){
	protocolFlagSet := flag.NewFlagSet("remote-forward", flag.ExitOnError)
	localAddress := protocolFlagSet.String( "local-address", "", "Address accessible by master")
	localPort := protocolFlagSet.String("local-port", "", "Port of the address that is accessible by master")
	masterAddress := protocolFlagSet.String("master-address", "", "Address of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return localAddress, localPort, masterAddress, masterPort
}


func ShowGeneralHelpMessage(){
	_, _ = fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " PROTOCOL *FLAGS\nProtocols available:\n\t - socks5\n\t - http\n\t - local-forward\n\t - remote-forward\n\t - master\n\t - translate")
}

func ShowTranslateHelpMessage(){
	_, _ = fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " translate TARGET *FLAGS\nTARGETS available:\n\t - forward-socks5\n\t")
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


func ParseMasterArguments() (*string, *string, *string, *string){
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	address := protocolFlagSet.String( "address", "0.0.0.0", "Address to listen on.")
	port := protocolFlagSet.String( "port", "1080", "Port to listen on.")
	remoteAddress := protocolFlagSet.String( "remote-address", "", "Argument required to handle correctly the \"remote-forward\"")
	remotePort := protocolFlagSet.String( "remote-port", "", "Argument required to handle correctly the \"remote-forward\"")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return address, port, remoteAddress, remotePort
}

func ParseForwardToSocks5Arguments() (*string, *string, *string, *string, *string, *string, *string, *string){
	protocolFlagSet := flag.NewFlagSet("forward-socks5", flag.ExitOnError)
	bindAddress := protocolFlagSet.String( "bind-address", "0.0.0.0", "Address to listen on.")
	bindPort := protocolFlagSet.String( "bind-port", "8080", "Port to listen on.")
	socks5Address := protocolFlagSet.String( "socks5-address", "127.0.0.1", "SOCKS5 server address to use")
	socks5Port := protocolFlagSet.String( "socks5-port", "1080", "SOCKS5 server port to use")
	username := protocolFlagSet.String( "socks5-username", "", "Username for the SOCKS5 server; leave black to no AUTH")
	password := protocolFlagSet.String( "socks5-password", "", "Password for the SOCKS5 server; leave black to no AUTH")
	targetAddress := protocolFlagSet.String( "target-address", "", "Address of the target host that is accessible by the SOCKS5 proxy")
	targetPort := protocolFlagSet.String( "target-port", "", "Port of the target host that is accessible by the SOCKS5 proxy")
	_ = protocolFlagSet.Parse(os.Args[3:])
	return bindAddress, bindPort, socks5Address, socks5Port, username, password, targetAddress, targetPort
}
