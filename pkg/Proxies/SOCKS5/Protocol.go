package SOCKS5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"math/big"
	"net"
)

type Socks5 struct {
	AuthenticationMethod ConnectionControllers.AuthenticationMethod
	WantedAuthMethod     byte
}

func ReceiveTargetRequest(clientConnectionReader *bufio.Reader) (byte, byte, []byte, []byte) {
	numberOfBytesReceived, targetRequest, ConnectionError := Sockets.Receive(clientConnectionReader, 1024)
	if ConnectionError == nil {
		if targetRequest[0] == Version {
			if targetRequest[1] == Connect || targetRequest[1] == Bind || targetRequest[1] == UDPAssociate {
				if targetRequest[3] == IPv4 || targetRequest[3] == IPv6 || targetRequest[3] == DomainName {
					return targetRequest[1], targetRequest[3], targetRequest[4 : numberOfBytesReceived-2], targetRequest[numberOfBytesReceived-2 : numberOfBytesReceived]
				}
			}
		}
	}
	return 0, 0, nil, nil
}

func GetTargetHostPort(targetRequestedCommand *byte, targetHostType *byte, rawTargetHost []byte, rawTargetPort []byte) (byte, string, string) {
	if *targetRequestedCommand != 0 && *targetHostType != 0 {
		switch *targetHostType {
		case IPv4:
			return *targetRequestedCommand, net.IPv4(rawTargetHost[0], rawTargetHost[1], rawTargetHost[2], rawTargetHost[3]).String(), new(big.Int).SetBytes(rawTargetPort).String()
		case IPv6:
			return *targetRequestedCommand, Sockets.GetIPv6(rawTargetHost), new(big.Int).SetBytes(rawTargetPort).String()
		case DomainName:
			return *targetRequestedCommand, string(rawTargetHost[1:]), new(big.Int).SetBytes(rawTargetPort).String()
		}
	}
	return ConnectionRefused, "", ""
}

func (socks5 *Socks5) SetAuthenticationMethod(authenticationMethod ConnectionControllers.AuthenticationMethod) error {
	socks5.AuthenticationMethod = authenticationMethod
	return nil
}

func (socks5 *Socks5) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {

	var targetRequestedCommand byte

	// Receive connection
	clientHasCompatibleAuthMethods := socks5.GetClientAuthenticationImplementedMethods(clientConnectionReader, clientConnectionWriter)
	if clientHasCompatibleAuthMethods {
		var targetHost string
		var targetPort string
		rawTargetRequestedCommand, targetHostType, rawTargetHost, rawTargetPort := ReceiveTargetRequest(
			clientConnectionReader)
		targetRequestedCommand, targetHost, targetPort = GetTargetHostPort(
			&rawTargetRequestedCommand, &targetHostType,
			rawTargetHost, rawTargetPort)
		if targetRequestedCommand != ConnectionRefused {
			return socks5.HandleCommandExecution(clientConnection, clientConnectionReader, clientConnectionWriter,
				&targetRequestedCommand, &targetHostType, &targetHost, &targetPort)
		}
	}
	var finalError string
	if !clientHasCompatibleAuthMethods {
		finalError = "No compatible auth methods found"
		_ = clientConnection.Close()
	} else if targetRequestedCommand == ConnectionRefused {
		finalError = "connection refused to target host"
	}
	return errors.New(finalError)
}
