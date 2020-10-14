package SOCKS5

import (
	"bufio"
	"net"
)

func HandleCommandExecution(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetRequestedCommand *byte, targetAddressType *byte,
	targetAddress *string, targetPort *string) {

	switch *targetRequestedCommand {
	case Connect:
		PrepareConnect(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	case Bind:
		PrepareBind(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	case UDPAssociate:
		PrepareUDPAssociate(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, targetAddressType)
	}
}
