package main

import (
	"FullProxy/FullProxy/ArgumentParser"
	"FullProxy/FullProxy/Interface"
	"FullProxy/FullProxy/Proxies/SOCKS5"
	"fmt"
)


func main(){
	arguments := ArgumentParser.GetArguments()
	if arguments.Parsed {
		switch arguments.Protocol {
		case "socks5":
			SOCKS5.StartSocks5(arguments.IP,
				arguments.Port,
				arguments.InterfaceMode,
				[]byte(arguments.Username),
				[]byte(arguments.Password))
		case "http":
			fmt.Println("HTTP not implemented yet")
		case "interface-master":
			Interface.Server(arguments.IP, arguments.Port)
		}
	}
}
