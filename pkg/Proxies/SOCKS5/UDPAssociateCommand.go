package SOCKS5

import (
	"net"
)

func PrepareUDPAssociate(
	clientConnection net.Conn, clientConnectionReader interface{},
	clientConnectionWriter interface{}, targetAddress *string,
	targetPort *string, targetAddressType *byte) {
	_ = clientConnection.Close()
}
