package ArgumentParser

import (
	"flag"
)


func GetArguments() map[string]interface{} {
	protocol := flag.String("protocol", "help", "Available protocols (socks4 socks5 http interface)")
	ip := flag.String("ip", "0.0.0.0", "IP to connect in interface mode or IP to bind in bind-mode. Default is 0.0.0.0")
	port := flag.String("port", "9050", "PORT to connect in interface mode or PORT to bind in bind-mode. Default is 9050")
	username := flag.String("username", "", "Username to force the client to authenticate with, default is not set")
	password := flag.String("password", "", "Password to force  the client to authenticate with, default is not set")
	interfaceMode := flag.Bool("interface-mode", false, "Use the interface mode")
	flag.Parse()

	var arguments = map[string]interface{}{
		"protocol":       *protocol,
		"ip":             *ip,
		"port":           *port,
		"interface-mode": *interfaceMode,
		"username": *username,
		"password": *password,
	}
	return arguments
}
