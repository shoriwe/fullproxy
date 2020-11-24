package SOCKS5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"net"
)

func (socks5 *Socks5) PrepareBind(
	clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer, targetHost *string,
	targetPort *string, targetHostType *byte) error {
	_ = clientConnection.Close()
	Templates.LogData(socks5.LoggingMethod, "Bind method not implemented yet")
	return errors.New("Bind method not implemented yet")
}
