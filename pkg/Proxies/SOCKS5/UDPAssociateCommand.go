package SOCKS5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"net"
)

func (socks5 *Socks5) PrepareUDPAssociate(clientConnection net.Conn, clientConnectionReader *bufio.Reader, clientConnectionWriter *bufio.Writer, targetHost *string,
	targetPort *string, targetHostType *byte) error {
	_ = clientConnection.Close()
	ConnectionControllers.LogData(socks5.LoggingMethod, "UDP-Associate method not implemented yet")
	return errors.New("UDP-Associate method not implemented yet")
}
