package circuit

import (
	"net"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type Conn struct {
	CloseFunctions []network.CloseFunc
	net.Conn
}

func (c *Conn) Close() (err error) {
	for _, closeFunc := range c.CloseFunctions {
		closeFunc()
	}
	if c.Conn != nil {
		err = c.Conn.Close()
	}
	return err
}
