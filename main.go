package main

import (
	"fmt"
	"github.com/shoriwe/FullProxy/src/ArgumentParser"
	"github.com/shoriwe/FullProxy/src/Interface"
	"github.com/shoriwe/FullProxy/src/Proxies/SOCKS5"
	"os"
)


func main() {
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintln(os.Stderr, "Try:\n", os.Args[0], " help")
		os.Exit(-1)
	}
	switch os.Args[1] {
	case "socks5":
		ip, port, username, password, interfaceMode := ArgumentParser.ParseSocks5Arguments()
		SOCKS5.StartSocks5(ip, port, interfaceMode, username, password)
	case "http":
		fmt.Println("HTTP not implemented yet")
	case "forward":

	case "interface-master":
		ip, port := ArgumentParser.ParseInterfaceArguments()
		Interface.Server(ip, port)
	case "help":
		ArgumentParser.ShowGeneralHelpMessage()
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Unknown command\nTry: ", os.Args[0], " help")
		os.Exit(-1)
	}
}
