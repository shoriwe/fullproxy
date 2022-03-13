package socks5

import (
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
)

func (socks5 *Socks5) UDPAssociate(sessionChunk []byte, clientConnection net.Conn) error {
	_ = clientConnection.Close()
	global.LogData(socks5.LoggingMethod, "UDP-Associate method not implemented yet")
	return errors.New("UDP-Associate method not implemented yet")
}
