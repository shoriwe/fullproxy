package SOCKS5

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"net"
)


func HandleCommandExecution(
	clientConnection net.Conn,
	clientConnectionReader ConnectionStructures.SocketReader, clientConnectionWriter ConnectionStructures.SocketWriter,
	targetRequestedCommand *byte, targetAddressType *byte,
	targetAddress *string, targetPort *string){


	switch *targetRequestedCommand {
	case Connect:
		PrepareConnect(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	case Bind:
		PrepareBind(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	case UDPAssociate:
		PrepareUDPAssociate(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	}
}
