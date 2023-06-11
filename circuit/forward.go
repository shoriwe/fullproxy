package circuit

import (
	"net"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type Forward struct {
	Network string
	Address string
}

func (f *Forward) Next(dial network.DialFunc) (closeFunc network.CloseFunc, newDial network.DialFunc, err error) {
	closeFunc = network.NopClose
	newDial = func(n, a string) (net.Conn, error) {
		return dial(f.Network, f.Address)
	}
	return closeFunc, newDial, err
}

// newForward ensures compile time safety, should never be used
func newForward() Knot {
	return &Forward{}
}
