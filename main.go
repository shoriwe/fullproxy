package main

import (
	"FullProxy/FullProxy/ArgumentParser"
	"FullProxy/FullProxy/Proxies/SOCKS5"
	"fmt"
	"os"
)


func main(){
	var arguments = ArgumentParser.GetArguments()
	switch arguments["protocol"].(string) {
	case "socks4":
		fmt.Println("Sock4 not implemented yet")
	case "socks5":
		SOCKS5.StartSocks5(arguments["ip"].(string),
							arguments["port"].(string),
							arguments["interface-mode"].(bool),
							[]byte(arguments["username"].(string)),
							[]byte(arguments["password"].(string)))
	case "http":
		fmt.Println("HTTP not implemented yet")
	case "interface":
		fmt.Println("Interface not implemented yet")
	default:
		_, _ = fmt.Fprint(os.Stderr, "Unknown module supplied, use \"help\"")
		os.Exit(1)
	}
}
