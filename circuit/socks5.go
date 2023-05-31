package circuit

import (
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"golang.org/x/net/proxy"
)

type Socks5 struct {
	Network string
	Address string
	Auth    *proxy.Auth
}

func (s *Socks5) Next(dial network.DialFunc) (closeFunc network.CloseFunc, newDial network.DialFunc, err error) {
	var dialer proxy.Dialer
	dialer, err = proxy.SOCKS5(s.Network, s.Address, s.Auth, &Dialer{DialFunc: dial})
	if err == nil {
		closeFunc = network.NopClose
		newDial = dialer.Dial
	}
	return closeFunc, newDial, err
}

// newSocks5 ensures compile time safety, should never be used
func newSocks5() Knot {
	return &Socks5{}
}
