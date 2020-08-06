package SOCKS5

import (
	"fmt"
	"net"
)


func PrepareUDPAssociate(
	clientConnection net.Conn, clientConnectionReader interface{},
	clientConnectionWriter interface{}, targetAddress *string,
	targetPort *string, targetAddressType *byte){
	_ = clientConnection.Close()
	fmt.Println(*targetAddress, *targetPort)
}
