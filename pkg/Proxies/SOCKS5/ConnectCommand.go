package SOCKS5

import (
	"bufio"
	"encoding/binary"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

func (socks5 *Socks5) PrepareConnect(
	clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetHost *string,
	targetPort *string,
	targetHostType *byte) error {

	targetConnection, connectionError := Sockets.Connect(targetHost, targetPort) // new(big.Int).SetBytes(rawTargetPort).String())
	if connectionError == nil {
		localHostPort := targetConnection.LocalAddr().(*net.TCPAddr)
		localPort := make([]byte, 2)
		binary.BigEndian.PutUint16(localPort, uint16(localHostPort.Port))
		response := []byte{Version, Succeeded, 0, *targetHostType}
		response = append(response[:], localHostPort.IP[:]...)
		response = append(response[:], localPort[:]...)
		_, connectionError = Sockets.Send(clientConnectionWriter, &response)
		if connectionError == nil {
			targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
			portProxy := Basic.PortProxy{
				TargetConnection:       targetConnection,
				TargetConnectionReader: targetConnectionReader,
				TargetConnectionWriter: targetConnectionWriter,
			}
			return portProxy.Handle(
				clientConnection,
				clientConnectionReader,
				clientConnectionWriter)
		} else {
			_ = clientConnection.Close()
			_ = targetConnection.Close()
			return connectionError
		}
	}
	failResponse := []byte{Version, ConnectionRefused, 0, *targetHostType, 0, 0}
	_, _ = Sockets.Send(clientConnectionWriter, &failResponse)
	_ = clientConnection.Close()
	return connectionError
}
