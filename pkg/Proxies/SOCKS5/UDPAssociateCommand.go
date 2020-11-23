package SOCKS5

import (
	"bufio"
	"errors"
	"net"
)

func (socks5 *Socks5) PrepareUDPAssociate(clientConnection net.Conn, clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer, targetHost *string,
	targetPort *string, targetHostType *byte) error {
	_ = clientConnection.Close()
	return errors.New("method not implemented yet")
}
