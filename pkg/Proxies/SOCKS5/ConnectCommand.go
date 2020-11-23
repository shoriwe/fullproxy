package SOCKS5

import (
	"bufio"
	"encoding/binary"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
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
	if connectionError != nil {
		ConnectionControllers.LogData(socks5.LoggingMethod, connectionError)
		failResponse := []byte{Version, ConnectionRefused, 0, *targetHostType, 0, 0}
		_, _ = Sockets.Send(clientConnectionWriter, &failResponse)
		_ = clientConnection.Close()
		return connectionError
	}
	localAddress := targetConnection.LocalAddr().(*net.TCPAddr)
	localPort := make([]byte, 2)
	binary.BigEndian.PutUint16(localPort, uint16(localAddress.Port))
	response := []byte{Version, Succeeded, 0, *targetHostType}
	response = append(response[:], localAddress.IP[:]...)
	response = append(response[:], localPort[:]...)
	_, connectionError = Sockets.Send(clientConnectionWriter, &response)
	if connectionError != nil {
		_ = clientConnection.Close()
		_ = targetConnection.Close()
		ConnectionControllers.LogData(socks5.LoggingMethod, connectionError)
		return connectionError
	}
	ConnectionControllers.LogData(socks5.LoggingMethod, "Client: ", clientConnection.RemoteAddr().String(), "  -> Target: ", targetConnection.RemoteAddr().String())
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	rawProxy := RawProxy.RawProxy{
		TargetConnection:       targetConnection,
		TargetConnectionReader: targetConnectionReader,
		TargetConnectionWriter: targetConnectionWriter,
		Tries:                  ConnectionControllers.GetTries(socks5.Tries),
		Timeout:                ConnectionControllers.GetTimeout(socks5.Timeout),
	}
	return rawProxy.Handle(
		clientConnection,
		clientConnectionReader,
		clientConnectionWriter)
}
