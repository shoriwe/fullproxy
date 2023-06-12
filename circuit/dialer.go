package circuit

import (
	"net"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type Dialer struct {
	DialFunc network.DialFunc
}

func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialFunc(network, address)
}
