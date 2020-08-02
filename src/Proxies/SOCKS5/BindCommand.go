package SOCKS5

import (
	"bufio"
	"net"
)


func PrepareBind(clientConnection net.Conn, clientConnectionReader *bufio.Reader,
clientConnectionWriter *bufio.Writer, targetAddress *string,
targetPort *string, targetAddressType *byte){
	_ = clientConnection.Close()
}