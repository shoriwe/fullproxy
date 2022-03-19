package socks5

import (
	"encoding/binary"
	"net"
)

func (socks5 *Socks5) UDPAssociate(sessionChunk []byte, clientConnection net.Conn, port int, host, hostPort string) error {
	clientRelayAddress, resolveError := net.ResolveUDPAddr("udp", hostPort)
	if resolveError != nil {
		_, _ = clientConnection.Write([]byte{SocksV5, GeneralSocksServerFailure, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		return resolveError
	}
	clientRelay, connectionError := net.DialUDP("udp",
		&net.UDPAddr{
			IP:   socks5.UDPRelay.LocalAddr().(*net.UDPAddr).IP,
			Port: DefaultRelayPort, // TODO: Maybe change this to a dynamic port selection
		},
		clientRelayAddress,
	)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{SocksV5, GeneralSocksServerFailure, 0x00, IPv4, 127, 0, 0, 1, 90, 90})
		return connectionError
	}

	var (
		addressType byte = 0x00 // FIXME
		localPort   [2]byte
	)
	bndAddress := socks5.UDPRelay.LocalAddr().(*net.UDPAddr).IP
	binary.BigEndian.PutUint16(localPort[:], DefaultRelayPort)

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
	clientRelay.Close()
	return nil
}
