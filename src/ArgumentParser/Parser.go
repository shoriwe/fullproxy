package ArgumentParser


import (
	"flag"
	"fmt"
	"os"
)


func PortForwardArguments() {

}


func ShowGeneralHelpMessage(){
	_, _ = fmt.Fprintln(os.Stderr, "Usage: ", os.Args[0], " PROTOCOL *FLAGS\nProtocols available:\n\t - socks5\n\t - http\n\t - interface-master")
}


func ParseSocks5Arguments() (string, string, []byte, []byte, bool) {
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	ip := *protocolFlagSet.String( "ip", "", "IP address to listen on. When \"-interface-mode\" flag is set, is the IP of interface to connect")
	port := *protocolFlagSet.String("port", "1080", "Port address to listen on. When \"-interface-mode\" flag is set, is the Port of the interface to connect. I both modes the default port is 1080")
	username := *protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := *protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	interfaceMode := *protocolFlagSet.Bool("interface-mode", false, "Connect to an interface, no bind proxying")
	if len(ip) == 0{
		if interfaceMode{
			ip = "127.0.0.1"
		} else {
			ip = "0.0.0.0"
		}
	}
	return ip, port, []byte(username), []byte(password), interfaceMode
}


func ParseInterfaceArguments() (string, string){
	protocolFlagSet := flag.NewFlagSet("interface", flag.ExitOnError)
	ip := *protocolFlagSet.String( "ip", "0.0.0.0", "IP address to listen on.")
	port := *protocolFlagSet.String( "port", "1080", "Port address to listen on.")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return ip, port
}

