package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
)

type UserArguments struct {
	Protocol      string
	IP            string
	Port          string
	Username      string
	Password      string
	InterfaceMode bool
	Parsed		  bool
}


func showGeneralHelpMessage(){
	_, _ = fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " PROTOCOL *FLAGS\nProtocols available:\n\t - socks5\n\t - http\n\t - interface-master")
}



func parseSocks5Arguments() UserArguments {
	socks5Arguments := UserArguments{}
	socks5Arguments.Protocol = "socks5"
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	protocolFlagSet.StringVar(&socks5Arguments.IP, "ip", "0.0.0.0", "IP address to listen on. When \"-interface-mode\" flag is set, is the IP of interface to connect")
	protocolFlagSet.StringVar(&socks5Arguments.Port, "port", "1080", "Port address to listen on. When \"-interface-mode\" flag is set, is the Port of the interface to connect")
	protocolFlagSet.StringVar(&socks5Arguments.Username, "username", "", "Username of the running proxy, requires \"-password\" and can't be an empty string ('')")
	protocolFlagSet.StringVar(&socks5Arguments.Password, "password", "", "Password of the running proxy, requires \"-username\" and can't be an empty string ('')")
	protocolFlagSet.BoolVar(&socks5Arguments.InterfaceMode, "interface-mode", false, "Connect to an interface, no bind proxying")
	_ = protocolFlagSet.Parse(os.Args[2:])
	socks5Arguments.Parsed = true
	return socks5Arguments
}


func parseInterfaceArguments() UserArguments{
	interfaceArguments := UserArguments{Username: "", Password: ""}
	interfaceArguments.Protocol = "interface-master"
	protocolFlagSet := flag.NewFlagSet("interface", flag.ExitOnError)
	protocolFlagSet.StringVar(&interfaceArguments.IP, "ip", "0.0.0.0", "IP address to listen on.")
	protocolFlagSet.StringVar(&interfaceArguments.Port, "port", "1080", "Port address to listen on.")
	_ = protocolFlagSet.Parse(os.Args[2:])
	interfaceArguments.Parsed = true
	return interfaceArguments
}


func GetArguments() UserArguments {
	var arguments = UserArguments{}
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "socks5":
			arguments = parseSocks5Arguments()
		case "interface-master":
			arguments = parseInterfaceArguments()
		case "http":
			arguments.Parsed = true
		default:
			arguments.Parsed = false
			showGeneralHelpMessage()
		}
	} else {
		arguments.Parsed = false
		showGeneralHelpMessage()
	}
	return arguments
}
