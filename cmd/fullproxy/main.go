package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/shoriwe/FullProxy/internal/Tools"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse/Master"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse/Slave"
	"github.com/shoriwe/FullProxy/pkg/Proxies/HTTP"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"github.com/shoriwe/FullProxy/pkg/Proxies/SOCKS5"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Translation/ForwardToSocks5"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	c2Address = "127.0.0.1:9051"
)

func init() {
	if os.Getenv("C2Address") != "" {
		c2Address = os.Getenv("C2Address")
	}
}

func loadList(filePath string) map[string]uint8 {
	file, openError := os.Open(filePath)
	if openError != nil {
		panic(openError)
	}
	content, readError := io.ReadAll(file)
	if readError != nil {
		panic(readError)
	}
	contentAsString := strings.ReplaceAll(string(content), "\r", "")
	result := map[string]uint8{}
	for _, line := range strings.Split(contentAsString, "\n") {
		result[line] = 0
	}
	return result
}

func configInboundFilter(whiteList, blackList string) Types.IOFilter {
	if whiteList != "" {
		reference := loadList(whiteList)
		return func(host string) bool {
			_, found := reference[host]
			return found
		}
	} else if blackList != "" {
		reference := loadList(whiteList)
		return func(host string) bool {
			_, found := reference[host]
			return !found
		}
	}
	return nil
}

func configOutboundFilter(whiteList, blackList string) Types.IOFilter {
	if whiteList != "" {
		reference := loadList(whiteList)
		return func(host string) bool {
			_, found := reference[host]
			return found
		}
	} else if blackList != "" {
		reference := loadList(whiteList)
		return func(host string) bool {
			_, found := reference[host]
			return !found
		}
	}
	return nil
}

func configAuthMethod(command, usersFile string) Types.AuthenticationMethod {
	if command != "" {
		return func(username []byte, password []byte) (bool, error) {
			cmd := exec.Command(command, hex.EncodeToString(username), hex.EncodeToString(password))

			if err := cmd.Start(); err != nil {
				log.Fatalf("cmd.Start: %v", err)
			}

			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					// The program has exited with an exit code != 0

					// This works on both Unix and Windows. Although package
					// syscall is generally platform dependent, WaitStatus is
					// defined for both Unix and Windows and in both cases has
					// an ExitStatus() method with the same signature.
					if status, ok2 := exitErr.Sys().(syscall.WaitStatus); ok2 {
						log.Printf("Exit Status: %d", status.ExitStatus())
					}
					return false, nil
				} else {
					log.Fatalf("cmd.Wait: %v", err)
				}
			}
			return true, nil
		}
	} else if usersFile != "" {
		reference := Templates.LoadUsers(usersFile)
		return func(username []byte, password []byte) (bool, error) {
			passwordHash, found := reference[string(username)]
			if !found {
				return false, nil
			}
			return Templates.SHA3512(password) == passwordHash, nil
		}

	}
	return nil
}

func configSocks5() (Types.IOFilter, Types.ProxyProtocol, error) {
	if len(os.Args) < 5 {
		return nil, SOCKS5.NewSocks5(nil, log.Println, nil), nil
	}
	flagSet := flag.NewFlagSet("socks5", flag.ExitOnError)
	authCommand := flagSet.String("auth-cmd", "", "shell command to pass the hex encoded username and password, exit code 0 means login success")
	usersFiles := flagSet.String("users-file", "", "json file with username as keys and sha3-513 of the password as values")

	inboundBlackList := flagSet.String("inbound-blacklist", "", "plain text file list with all the HOST that are forbidden to connect to the proxy")
	inboundWhiteList := flagSet.String("inbound-whitelist", "", "plain text file list with all the HOST that are permitted to connect to the proxy")
	outboundBlackList := flagSet.String("outbound-blacklist", "", "plain text file list with all the forbidden proxy targets")
	outboundWhiteList := flagSet.String("outbound-whitelist", "", "plain text file list with all the permitted proxy targets")

	parsingError := flagSet.Parse(os.Args[5:])

	if parsingError != nil {
		panic(parsingError)
	}

	return configInboundFilter(*inboundWhiteList, *inboundBlackList), SOCKS5.NewSocks5(configAuthMethod(*authCommand, *usersFiles), log.Println, configOutboundFilter(*outboundWhiteList, *outboundBlackList)), nil
}

func configPortForward() (Types.IOFilter, Types.ProxyProtocol, error) {
	if len(os.Args) < 5 {
		return nil, SOCKS5.NewSocks5(nil, log.Println, nil), nil
	}
	flagSet := flag.NewFlagSet("port-forward", flag.ExitOnError)

	networkType := flagSet.String("network-type", "tcp", "tcp or udp")
	targetAddress := flagSet.String("target-address", "127.0.0.1:80", "Address to connect")

	inboundBlackList := flagSet.String("inbound-blacklist", "", "plain text file list with all the HOST that are forbidden to connect to the proxy")
	inboundWhiteList := flagSet.String("inbound-whitelist", "", "plain text file list with all the HOST that are permitted to connect to the proxy")

	parsingError := flagSet.Parse(os.Args[5:])

	if parsingError != nil {
		panic(parsingError)
	}

	return configInboundFilter(*inboundWhiteList, *inboundBlackList), PortForward.NewForward(*networkType, *targetAddress, log.Println), nil
}

func configHTTP() (Types.IOFilter, Types.ProxyProtocol, error) {
	if len(os.Args) < 5 {
		return nil, SOCKS5.NewSocks5(nil, log.Println, nil), nil
	}
	flagSet := flag.NewFlagSet("http", flag.ExitOnError)
	authCommand := flagSet.String("auth-cmd", "", "shell command to pass the hex encoded username and password, exit code 0 means login success")
	usersFiles := flagSet.String("users-file", "", "json file with username as keys and sha3-513 of the password as values")

	inboundBlackList := flagSet.String("inbound-blacklist", "", "plain text file list with all the HOST that are forbidden to connect to the proxy")
	inboundWhiteList := flagSet.String("inbound-whitelist", "", "plain text file list with all the HOST that are permitted to connect to the proxy")
	outboundBlackList := flagSet.String("outbound-blacklist", "", "plain text file list with all the forbidden proxy targets")
	outboundWhiteList := flagSet.String("outbound-whitelist", "", "plain text file list with all the permitted proxy targets")

	parsingError := flagSet.Parse(os.Args[5:])

	if parsingError != nil {
		panic(parsingError)
	}

	return configInboundFilter(*inboundWhiteList, *inboundBlackList), HTTP.NewHTTP(
		configAuthMethod(*authCommand, *usersFiles),
		log.Println,
		configOutboundFilter(*outboundWhiteList, *outboundBlackList),
	), nil
}

func configTranslateSocks5() (Types.IOFilter, Types.ProxyProtocol, error) {
	if len(os.Args) < 5 {
		return nil, SOCKS5.NewSocks5(nil, log.Println, nil), nil
	}
	flagSet := flag.NewFlagSet("translate-sock5", flag.ExitOnError)

	networkType := flagSet.String("network-type", "tcp", "tcp or udp")
	targetAddress := flagSet.String("target-address", "127.0.0.1:80", "Address to connect")

	socks5ProxyAddress := flagSet.String("socks5", "127.0.0.1:9050", "Address of the socks5 url")
	socks5Username := flagSet.String("username", "", "Username for the socks5 server")
	socks5Password := flagSet.String("password", "", "Password for the socks5 server")

	inboundBlackList := flagSet.String("inbound-blacklist", "", "plain text file list with all the HOST that are forbidden to connect to the proxy")
	inboundWhiteList := flagSet.String("inbound-whitelist", "", "plain text file list with all the HOST that are permitted to connect to the proxy")

	parsingError := flagSet.Parse(os.Args[5:])

	if parsingError != nil {
		panic(parsingError)
	}
	translate, generationError := ForwardToSocks5.NewForwardToSocks5(*networkType, *socks5ProxyAddress, *socks5Username, *socks5Password, *targetAddress, log.Println)
	if generationError != nil {
		return nil, nil, generationError
	}
	return configInboundFilter(*inboundWhiteList, *inboundBlackList), translate, nil
}

func configProtocol(protocol string) (Types.IOFilter, Types.ProxyProtocol, error) {
	switch protocol {
	case "socks5":
		return configSocks5()
	case "port-forward":
		return configPortForward()
	case "http":
		return configHTTP()
	case "translate-socks5":
		return configTranslateSocks5()
	}
	return nil, nil, errors.New("Unknown protocol: " + protocol)
}

func main() {
	numberOfArguments := len(os.Args)
	if numberOfArguments < 5 {
		_, _ = fmt.Fprintf(os.Stderr, "%s %s", os.Args[0],
			"MODE NETWORK_TYPE ADDRESS PROTOCOL [OPTIONS]\n"+
				"\t- MODE:         bind|master|slave\n"+
				"\t- NETWORK_TYPE: tcp|udp\n"+
				"\t- ADDRESS:      IPv4|IPv6 or Domain followed by \":\" and the PORT; For Example -> \"127.0.0.1:80\"\n"+
				"\t- PROTOCOL:     socks5|http|port-forward|translate-socks5\n"+
				"Environment Variables:\n"+
				"\t- C2Address     Host and port of the C2 port of the master server\n",
		)

		os.Exit(1)
	}

	mode := os.Args[1]
	networkType := os.Args[2]
	address := os.Args[3]
	protocol := os.Args[4]

	inboundFilter, proxyProtocol, setupError := configProtocol(protocol)
	if setupError != nil {
		log.Fatal(setupError)
	}
	var (
		pipe Types.Pipe
	)
	switch mode {
	case "bind":
		pipe = Pipes.NewBindPipe(networkType, address, proxyProtocol, log.Println, inboundFilter)
	case "master":
		pipe = Master.NewMaster(networkType, c2Address, address, log.Println, inboundFilter, proxyProtocol)
	case "slave":
		pipe = Slave.NewSlave(networkType, c2Address, log.Println)
	default:
		panic("Unknown mode")
	}
	log.Fatal(pipe.Serve())
}
