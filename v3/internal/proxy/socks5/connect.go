package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"net"
)

func (socks5 *Socks5) Connect(clientConnection net.Conn, _ int, host, hostPort string) error {
	if !global.FilterOutbound(socks5.OutboundFilter, host) {
		_, connectionError := clientConnection.Write([]byte{SocksV5, ConnectionNotAllowedByRuleSet, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		global.LogData(socks5.LoggingMethod, "Forbidden connection to: "+host)
		return connectionError
	}

	// Try to connect to the target

	targetConnection, connectionError := socks5.Dial("tcp", hostPort)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{SocksV5, GeneralSocksServerFailure, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		return connectionError
	}

	// Respond to client

	var (
		addressType byte = 0x00 // FIXME
		localPort   [2]byte
	)
	bndAddress := targetConnection.LocalAddr().(*net.TCPAddr).IP
	binary.BigEndian.PutUint16(localPort[:], uint16(targetConnection.LocalAddr().(*net.TCPAddr).Port))

	if bndAddress.To4() != nil {
		bndAddress = bndAddress.To4()
		addressType = IPv4
	} else if bndAddress.To16() != nil {
		bndAddress = bndAddress.To16()
		addressType = IPv6
	}

	response := make([]byte, 0, 350)
	response = append(response, SocksV5, ConnectionSucceed, 0x00, addressType)
	response = append(response, bndAddress...)
	response = append(response, localPort[:]...)

	_, connectionError = clientConnection.Write(response)
	if connectionError != nil {
		return connectionError
	}
	// Forward traffic
	return pipes.ForwardTraffic(clientConnection, targetConnection)
}
