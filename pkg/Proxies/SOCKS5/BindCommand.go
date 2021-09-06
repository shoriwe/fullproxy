package SOCKS5

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"net"
)

func (socks5 *Socks5) Bind(clientConnection net.Conn) error {
	_ = clientConnection.Close()
	Tools.LogData(socks5.LoggingMethod, "Bind method not implemented yet")
	return errors.New("Bind method not implemented yet")
}
