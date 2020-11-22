package SOCKS5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

func (socks5 *Socks5)PrepareConnect(
	targetAddress *string,
	targetPort *string,
	targetAddressType *byte) error {

	targetConnection, connectionError := Sockets.Connect(targetAddress, targetPort) // new(big.Int).SetBytes(rawTargetPort).String())
	if connectionError == nil {
		localAddressPort := targetConnection.LocalAddr().(*net.TCPAddr)
		localPort := make([]byte, 2)
		binary.BigEndian.PutUint16(localPort, uint16(localAddressPort.Port))
		response := []byte{Version, Succeeded, 0, *targetAddressType}
		response = append(response[:], localAddressPort.IP[:]...)
		response = append(response[:], localPort[:]...)
		_, connectionError = Sockets.Send(socks5.ClientConnectionWriter, &response)
		if connectionError == nil {
			targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
			portProxy := Basic.PortProxy{
				TargetConnection: targetConnection,
				TargetConnectionReader: targetConnectionReader,
				TargetConnectionWriter: targetConnectionWriter,
			}
			return portProxy.Handle(
				socks5.ClientConnection,
				socks5.ClientConnectionReader,
				socks5.ClientConnectionWriter)
		} else {
			_ = socks5.ClientConnection.Close()
			_ = targetConnection.Close()
			return connectionError
		}
	}
	failResponse := []byte{Version, ConnectionRefused, 0, *targetAddressType, 0, 0}
	_, _ = Sockets.Send(socks5.ClientConnectionWriter, &failResponse)
	_ = socks5.ClientConnection.Close()
	return connectionError
}
