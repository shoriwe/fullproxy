package SOCKS5

import (
	"errors"
	"github.com/shoriwe/FullProxy/v2/pkg/Tools"
	"net"
)

func (socks5 *Socks5) UDPAssociate(clientConnection net.Conn) error {
	_ = clientConnection.Close()
	Tools.LogData(socks5.LoggingMethod, "UDP-Associate method not implemented yet")
	return errors.New("UDP-Associate method not implemented yet")
}
