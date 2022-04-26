package socks5

import (
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
)

func (socks5 *Socks5) Bind(context *Context) error {
	_ = context.ClientConnection.Close()
	global.LogData(socks5.LoggingMethod, "Bind method not implemented yet")
	return errors.New("bind method not implemented yet")
}
