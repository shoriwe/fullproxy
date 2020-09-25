package SOCKS5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
)

func PrepareConnect(
	clientConnection net.Conn, clientConnectionReader ConnectionStructures.SocketReader,
	clientConnectionWriter ConnectionStructures.SocketWriter, targetAddress *string,
	targetPort *string, targetAddressType *byte) {

	targetConnection := Sockets.Connect(targetAddress, targetPort) // new(big.Int).SetBytes(rawTargetPort).String())
	if targetConnection != nil {
		localAddressPort := targetConnection.LocalAddr().(*net.TCPAddr)
		localPort := make([]byte, 2)
		binary.BigEndian.PutUint16(localPort, uint16(localAddressPort.Port))
		response := []byte{Version, Succeeded, 0, *targetAddressType}
		response = append(response[:], localAddressPort.IP[:]...)
		response = append(response[:], localPort[:]...)
		_, connectionError := Sockets.Send(clientConnectionWriter, &response)
		if connectionError == nil {
			targetConnectionReader, targetConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(targetConnection)
			Basic.Proxy(
				clientConnection, targetConnection,
				clientConnectionReader, clientConnectionWriter,
				targetConnectionReader, targetConnectionWriter)
		} else {
			{
				_ = clientConnection.Close()
				_ = targetConnection.Close()
			}
		}
	} else {
		failResponse := []byte{Version, ConnectionRefused, 0, *targetAddressType, 0, 0}
		_, _ = Sockets.Send(clientConnectionWriter, &failResponse)
		_ = clientConnection.Close()
	}
}
