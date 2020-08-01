package SOCKS5

import (
	"bufio"
	"net"
)


func PrepareUDPAssociate(clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer, targetAddress *string,
	targetPort *string, rawTargetAddress []byte,
	rawTargetPort []byte, targetAddressType *byte){
	_ = clientConnection.Close()
}