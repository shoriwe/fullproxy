package SOCKS5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"net"
	"strconv"
	"strings"
)

func (socks5 *Socks5) Connect(clientConnection net.Conn) error {
	reserved := make([]byte, 1)
	numberOfBytesReceived, connectionError := clientConnection.Read(reserved)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != 1 {
		return protocolError
	}
	addressType := make([]byte, 1)
	numberOfBytesReceived, connectionError = clientConnection.Read(addressType)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != 1 {
		return protocolError
	}

	var targetHostLength int
	switch addressType[0] {
	case IPv4:
		targetHostLength = 4
	case DomainName:
		domainLength := make([]byte, 1)
		numberOfBytesReceived, connectionError = clientConnection.Read(domainLength)
		if connectionError != nil {
			return connectionError
		} else if numberOfBytesReceived != 1 {
			return protocolError
		}
		targetHostLength = int(domainLength[0])
	case IPv6:
		targetHostLength = 16
	default:
		return protocolError
	}
	rawTargetHost := make([]byte, targetHostLength)
	numberOfBytesReceived, connectionError = clientConnection.Read(rawTargetHost)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != targetHostLength {
		return protocolError
	}

	rawTargetPort := make([]byte, 2)
	numberOfBytesReceived, connectionError = clientConnection.Read(rawTargetPort)
	if connectionError != nil {
		return connectionError
	} else if numberOfBytesReceived != 2 {
		return protocolError
	}

	// Cleanup the address

	target := clean(addressType[0], rawTargetHost, rawTargetPort)

	// Try to connect to the target

	var targetConnection net.Conn
	targetConnection, connectionError = net.Dial("tcp", target)
	if connectionError != nil {
		// Respond the error to the client
		return connectionError
	}

	// Respond to client

	local := net.ParseIP(targetConnection.LocalAddr().String()).To16()
	localAddressBytes, _ := local.MarshalText()
	response := []byte{SocksV5, Succeeded, 0, IPv6, byte(len(localAddressBytes))}
	response = append(response, localAddressBytes...)
	portAsInt, _ := strconv.Atoi(strings.ReplaceAll(targetConnection.LocalAddr().String(), local.String(), ""))
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(portAsInt))
	response = append(response, port...)
	_, connectionError = clientConnection.Write(response)
	if connectionError != nil {
		return connectionError
	}
	// Forward traffic
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}
