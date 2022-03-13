package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"net"
)

func (socks5 *Socks5) Connect(sessionChunk []byte, clientConnection net.Conn) error {
	var (
		targetHostLength, hostStartIndex int
	)
	switch sessionChunk[3] {
	case IPv4:
		targetHostLength = 4
		hostStartIndex = 4
	case DomainName:
		targetHostLength = int(sessionChunk[4])
		hostStartIndex = 5
	case IPv6:
		targetHostLength = 16
		hostStartIndex = 4
	default:
		return protocolError
	}
	rawTargetHost := sessionChunk[hostStartIndex : hostStartIndex+targetHostLength]

	rawTargetPort := sessionChunk[hostStartIndex+targetHostLength : hostStartIndex+targetHostLength+2]

	// Cleanup the address

	host, target := clean(sessionChunk[3], rawTargetHost, rawTargetPort)
	if !global.FilterOutbound(socks5.OutboundFilter, host) {
		_, connectionError := clientConnection.Write([]byte{SocksV5, ConnectionNotAllowedByRuleSet, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		global.LogData(socks5.LoggingMethod, "Forbidden connection to: "+host)
		return connectionError
	}

	// Try to connect to the target

	targetConnection, connectionError := socks5.Dial("tcp", target)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{SocksV5, GeneralSocksServerFailure, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		return connectionError
	}

	// Respond to client

	var (
		bndAddress net.IP
		localType  byte = 0x00 // FIXME
		localPort  [2]byte
	)
	if targetAsTCP, ok := targetConnection.LocalAddr().(*net.TCPAddr); ok {
		bndAddress = targetAsTCP.IP
		binary.BigEndian.PutUint16(localPort[:], uint16(targetAsTCP.Port))
	} else if targetAsUDP, ok := targetConnection.LocalAddr().(*net.TCPAddr); ok {
		bndAddress = targetAsUDP.IP
		binary.BigEndian.PutUint16(localPort[:], uint16(targetAsUDP.Port))
	} else {
		bndAddress = net.IPv4(127, 0, 0, 1)
		binary.BigEndian.PutUint16(localPort[:], 8080)
	}

	if bndAddress.To4() != nil {
		bndAddress = bndAddress.To4()
		localType = IPv4
	} else if bndAddress.To16() != nil {
		bndAddress = bndAddress.To4()
		localType = IPv6
	}
	if bndAddress == nil { // FIXME: This hack need to be removed
		bndAddress = net.IPv4(127, 0, 0, 1)
		binary.BigEndian.PutUint16(localPort[:], 8080)
		bndAddress = bndAddress.To4()
		localType = IPv4
	}

	response := make([]byte, 0, 350)
	response = append(response, SocksV5, ConnectionSucceed, 0x00, localType)
	response = append(response, bndAddress...)
	response = append(response, localPort[:]...)

	_, connectionError = clientConnection.Write(response)
	if connectionError != nil {
		return connectionError
	}
	// Forward traffic
	return pipes.ForwardTraffic(clientConnection, targetConnection)
}
