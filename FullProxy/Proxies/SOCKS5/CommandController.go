package SOCKS5

import (
	"bufio"
	"net"
)


func HandleCommandExecution(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer,
	targetRequestedCommand *byte, targetAddressType *byte,
	targetAddress *string, targetPort *string,
	rawTargetAddress []byte, rawTargetPort []byte) net.Conn{

	var targetConnection net.Conn = nil
	switch *targetRequestedCommand {
	case Connect:
		targetConnection = PrepareConnect(clientConnection, clientConnectionReader, clientConnectionWriter, targetAddress, targetPort, rawTargetAddress, rawTargetPort, targetAddressType)
	case Bind:
		break
	case UDPAssociate:
		break
	}
	return targetConnection
}
