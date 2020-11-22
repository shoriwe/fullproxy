package SOCKS5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"math/big"
	"net"
)

type Socks5 struct {
	AuthenticationMethod ConnectionHandlers.AuthenticationMethod
	ClientConnection net.Conn
	ClientConnectionReader *bufio.Reader
	ClientConnectionWriter *bufio.Writer
	WantedAuthMethod byte
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

func GetTargetAddressPort(targetRequestedCommand *byte, targetAddressType *byte, rawTargetAddress []byte, rawTargetPort []byte) (byte, string, string) {
	if *targetRequestedCommand != 0 && *targetAddressType != 0 {
		switch *targetAddressType {
		case IPv4:
			return *targetRequestedCommand, net.IPv4(rawTargetAddress[0], rawTargetAddress[1], rawTargetAddress[2], rawTargetAddress[3]).String(), new(big.Int).SetBytes(rawTargetPort).String()
		case IPv6:
			return *targetRequestedCommand, Sockets.GetIPv6(rawTargetAddress), new(big.Int).SetBytes(rawTargetPort).String()
		case DomainName:
			return *targetRequestedCommand, string(rawTargetAddress[1:]), new(big.Int).SetBytes(rawTargetPort).String()
		}
	}
	return ConnectionRefused, "", ""
}

func (socks5 *Socks5)Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {

	socks5.ClientConnection = clientConnection
	socks5.ClientConnectionReader = clientConnectionReader
	socks5.ClientConnectionWriter = clientConnectionWriter

	var targetRequestedCommand byte

	// Receive connection
	clientHasCompatibleAuthMethods := socks5.GetClientAuthenticationImplementedMethods()
	if clientHasCompatibleAuthMethods {
		var targetAddress string
		var targetPort string
		rawTargetRequestedCommand, targetAddressType, rawTargetAddress, rawTargetPort := ReceiveTargetRequest(
			clientConnectionReader)
		targetRequestedCommand, targetAddress, targetPort = GetTargetAddressPort(
			&rawTargetRequestedCommand, &targetAddressType,
			rawTargetAddress, rawTargetPort)
		if targetRequestedCommand != ConnectionRefused {
			return socks5.HandleCommandExecution(&targetRequestedCommand,
				&targetAddressType, &targetAddress, &targetPort)
		}
	}
	var finalError string
	if !clientHasCompatibleAuthMethods {
		finalError = "No compatible auth methods found"
		_ = clientConnection.Close()
	} else if targetRequestedCommand == ConnectionRefused {
		finalError = "connection refused to target address"
	}
	return errors.New(finalError)
}
