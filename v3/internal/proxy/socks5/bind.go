package socks5

import (
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
)

func (socks5 *Socks5) Bind(sessionChunk []byte, clientConnection net.Conn, port int, host, target string) error {
	_ = clientConnection.Close()
	global.LogData(socks5.LoggingMethod, "Bind method not implemented yet")
	return errors.New("Bind method not implemented yet")
}
