package SOCKS5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Tools"
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

	host, target := clean(addressType[0], rawTargetHost, rawTargetPort)
	if !Tools.FilterOutbound(socks5.OutboundFilter, host) {
		Tools.LogData(socks5.LoggingMethod, "Forbidden connection to: "+host)
		return nil
	}
	// Try to connect to the target

	var targetConnection net.Conn
	targetConnection, connectionError = net.Dial("tcp", target)
	if connectionError != nil {
		// Respond the error to the client
		return connectionError
	}

	// Respond to client

	local := strings.Split(targetConnection.LocalAddr().String(), ":")
	var localAddressBytes []byte
	for _, rawNumber := range strings.Split(local[0], ".") {
		number, _ := strconv.Atoi(rawNumber)
		localAddressBytes = append(localAddressBytes, uint8(number))
	}
	response := []byte{SocksV5, 0x00, 0, IPv4}
	response = append(response, localAddressBytes...)
	portAsInt, _ := strconv.Atoi(local[1])
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
