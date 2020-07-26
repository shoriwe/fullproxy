package SOCKS5

import (
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"net"
	"time"
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
	HandleConnect(
		clientConnection, targetConnection,
		clientConnectionReader, clientConnectionWriter,
		targetConnectionReader, targetConnectionWriter,
	)
	return targetConnection
}

func HandleConnect(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetConnectionReader *bufio.Reader,
	targetConnectionWriter *bufio.Writer) {
	for {
		for state := 0; state < 2; state++ {
			switch state {
			case 1:
				_ = clientConnection.SetReadDeadline(time.Now().Add(100))
				NumberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(clientConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				_, ConnectionError = Sockets.Send(targetConnectionWriter, buffer[:NumberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			case 0:
				_ = targetConnection.SetReadDeadline(time.Now().Add(100))
				NumberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(targetConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				_, ConnectionError = Sockets.Send(clientConnectionWriter, buffer[:NumberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			}
		}
	}
}
