package SOCKS5

import (
	"FullProxy/FullProxy/Proxies/Basic"
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"net"
)


func PrepareConnect(
	clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer, targetAddress *string,
	targetPort *string, rawTargetAddress []byte,
	rawTargetPort []byte, targetAddressType *byte) net.Conn{

	var targetConnection = Sockets.Connect(*targetAddress, *targetPort) // new(big.Int).SetBytes(rawTargetPort).String())
	if targetConnection == nil {
		_, _ = Sockets.Send(clientConnectionWriter, []byte{Version, ConnectionRefused, 0, *targetAddressType, 0, 0})
		return nil
	}
	response := []byte{Version, Succeeded, 0, *targetAddressType}
	response = append(response[:], rawTargetAddress[:]...)
	response = append(response[:], rawTargetPort[:]...)
	_, ConnectionError := Sockets.Send(clientConnectionWriter, response)
	if ConnectionError != nil {
		return nil
	}
	targetConnectionReader := bufio.NewReader(targetConnection)
	targetConnectionWriter := bufio.NewWriter(targetConnection)
	Basic.Proxy(
		clientConnection, targetConnection,
		clientConnectionReader, clientConnectionWriter,
		targetConnectionReader, targetConnectionWriter, true)
	return targetConnection
}
