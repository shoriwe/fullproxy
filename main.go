package main

import (
	"FullProxy/ArgumentParser"
	"FullProxy/Proxies"
	"fmt"
	"os"
)


func main(){
	var arguments = ArgumentParser.GetArguments()
	switch arguments["protocol"].(string) {
	case "socks4":
		fmt.Println("Sock4 not implemented yet")
	case "socks5":
		Proxies.StartSocks5(arguments["ip"].(string),
							arguments["port"].(string),
							arguments["interface-mode"].(bool),
							arguments["username"].(string),
							arguments["password"].(string))
	case "http":
		fmt.Println("HTTP not implemented yet")
	case "interface":
		fmt.Println("Interface not implemented yet")
	default:
		_, _ = fmt.Fprint(os.Stderr, "Unknown module supplied, use \"help\"")
		os.Exit(1)
	}
}
