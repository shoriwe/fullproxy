package SOCKS5

import (
	"bufio"
	"encoding/binary"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
)


func PrepareConnect(
	clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer, targetAddress *string,
	targetPort *string, targetAddressType *byte){

	var connectionError error
	var targetConnection net.Conn
	targetConnection = Sockets.Connect(*targetAddress, *targetPort) // new(big.Int).SetBytes(rawTargetPort).String())
	localAddressPort := targetConnection.LocalAddr().(*net.TCPAddr)
	localPort := make([]byte, 2)
	binary.BigEndian.PutUint16(localPort, uint16(localAddressPort.Port))
	if targetConnection != nil {
		response := []byte{Version, Succeeded, 0, *targetAddressType}
		response = append(response[:], localAddressPort.IP[:]...)
		response = append(response[:], localPort[:]...)
		_, connectionError = Sockets.Send(clientConnectionWriter, response)
		if connectionError == nil {
			targetConnectionReader := bufio.NewReader(targetConnection)
			targetConnectionWriter := bufio.NewWriter(targetConnection)
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
		_, _ = Sockets.Send(clientConnectionWriter, []byte{Version, ConnectionRefused, 0, *targetAddressType, 0, 0})
		_ = clientConnection.Close()
	}
}
