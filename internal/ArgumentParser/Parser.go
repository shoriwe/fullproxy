package ArgumentParser

import (
	"flag"
	"fmt"
	"os"
)

func HTTPArguments() (*string, *string, []byte, []byte, *bool, *bool) {
	protocolFlagSet := flag.NewFlagSet("http", flag.ExitOnError)
	host := protocolFlagSet.String("host", "", "Host to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "8080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 8080")
	username := protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	tls := protocolFlagSet.Bool("tls", false, "Use HTTPS")
	_ = protocolFlagSet.Parse(os.Args[2:])
	if len(*host) == 0 {
		if *slave {
			*host = "127.0.0.1"
		} else {
			*host = "0.0.0.0"
		}
	}
	return host, port, []byte(*username), []byte(*password), slave, tls
}

func LocalPortForwardArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("local-forward", flag.ExitOnError)
	forwardHost := protocolFlagSet.String("forward-host", "", "Host to forward the traffic received from master")
	forwardPort := protocolFlagSet.String("forward-port", "", "Port to forward the traffic received from master")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return forwardHost, forwardPort, masterHost, masterPort
}

func RemotePortForwardArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("remote-forward", flag.ExitOnError)
	localHost := protocolFlagSet.String("local-host", "", "Host to bind by slave")
	localPort := protocolFlagSet.String("local-port", "", "Port to bind by slave")
	masterHost := protocolFlagSet.String("master-host", "", "Host of the master")
	masterPort := protocolFlagSet.String("master-port", "", "Port of the master")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return localHost, localPort, masterHost, masterPort
}

func ShowGeneralHelpMessage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:\n\t", os.Args[0], "PROTOCOL *FLAGS\n\nProtocols available:\n\t - socks5\n\t - http\n\t - local-forward\n\t - remote-forward\n\t - master\n\t - translate")
}

func ShowTranslateHelpMessage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage:\n\t", os.Args[0], "translate TARGET *FLAGS\n\nTARGETS available:\n\t - port_forward-socks5\n\t")
}

func ParseSocks5Arguments() (*string, *string, []byte, []byte, *bool) {
	protocolFlagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	host := protocolFlagSet.String("host", "", "Host to listen on. When \"-slave\" flag is set, is the IP of master to connect")
	port := protocolFlagSet.String("port", "1080", "Port to listen on. When \"-slave\" flag is set, is the Port of the master to connect. I both modes the default port is 1080")
	username := protocolFlagSet.String("username", "", "Username of the running proxy, requires \"-password\". It will be ignored if is an empty string")
	password := protocolFlagSet.String("password", "", "Password of the running proxy, requires \"-username\". It will be ignored if is an empty string")
	slave := protocolFlagSet.Bool("slave", false, "Connect to a master, no bind proxying")
	_ = protocolFlagSet.Parse(os.Args[2:])
	if len(*host) == 0 {
		if *slave {
			*host = "127.0.0.1"
		} else {
			*host = "0.0.0.0"
		}
	}
	return host, port, []byte(*username), []byte(*password), slave
}

func ParseMasterArguments() (*string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("master", flag.ExitOnError)
	host := protocolFlagSet.String("host", "0.0.0.0", "Host to listen on.")
	port := protocolFlagSet.String("port", "1080", "Port to listen on.")
	remoteHost := protocolFlagSet.String("forward-host", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	remotePort := protocolFlagSet.String("forward-port", "", "Argument required to handle correctly the \"remote-forward\" (This is the service that the master can only acceded)")
	_ = protocolFlagSet.Parse(os.Args[2:])
	return host, port, remoteHost, remotePort
}

func ParseForwardToSocks5Arguments() (*string, *string, *string, *string, *string, *string, *string, *string) {
	protocolFlagSet := flag.NewFlagSet("port_forward-socks5", flag.ExitOnError)
	bindHost := protocolFlagSet.String("bind-host", "0.0.0.0", "Host to listen on.")
	bindPort := protocolFlagSet.String("bind-port", "8080", "Port to listen on.")
	socks5Host := protocolFlagSet.String("socks5-host", "127.0.0.1", "SOCKS5 server host to use")
	socks5Port := protocolFlagSet.String("socks5-port", "1080", "SOCKS5 server port to use")
	username := protocolFlagSet.String("socks5-username", "", "Username for the SOCKS5 server; leave empty for no AUTH")
	password := protocolFlagSet.String("socks5-password", "", "Password for the SOCKS5 server; leave empty for no AUTH")
	targetHost := protocolFlagSet.String("target-host", "", "Host of the target host that is accessible by the SOCKS5 proxy")
	targetPort := protocolFlagSet.String("target-port", "", "Port of the target host that is accessible by the SOCKS5 proxy")
	_ = protocolFlagSet.Parse(os.Args[3:])
	return bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort
}
