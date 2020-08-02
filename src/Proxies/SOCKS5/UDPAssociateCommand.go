package SOCKS5

import (
	"bufio"
	"fmt"
	"net"
)


func PrepareUDPAssociate(clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer, targetAddress *string,
	targetPort *string, targetAddressType *byte){
	_ = clientConnection.Close()
	fmt.Println(*targetAddress, *targetPort)
}
