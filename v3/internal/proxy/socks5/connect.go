package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"net"
	"strconv"
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
		global.LogData(socks5.LoggingMethod, "Forbidden connection to: "+host)
		return nil
	}

	// Try to connect to the target

	targetConnection, connectionError := socks5.Dial("tcp", target)
	if connectionError != nil {
		// TODO: Respond the error to the client
		return connectionError
	}

	// Respond to client

	var (
		h, p, _             = net.SplitHostPort(targetConnection.LocalAddr().String())
		numericLocalPort, _ = strconv.Atoi(p)
		bndAddress          = net.ParseIP(h)
		localType           byte
		localPort           [2]byte
	)
	binary.BigEndian.PutUint16(localPort[:], uint16(numericLocalPort))
	if bndAddress.To4() == nil {
		bndAddress = []byte(bndAddress)[len([]byte(bndAddress))-net.IPv6len:]
		localType = IPv6
	} else {
		bndAddress = []byte(bndAddress)[len([]byte(bndAddress))-net.IPv4len:]
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
