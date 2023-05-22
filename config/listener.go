package config

import (
	"crypto/tls"
	"net"
)

type Listener struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
	TLS     *TLS   `yaml:"tls"`
}

func (l *Listener) Listen() (net.Listener, error) {
	if l.TLS == nil {
		return net.Listen(l.Network, l.Address)
	}
	config, cErr := l.TLS.Config()
	if cErr != nil {
		return nil, cErr
	}
	return tls.Listen(l.Network, l.Address, config)
}

func (l *Listener) Dial() (net.Conn, error) {
	if l.TLS == nil {
		return net.Dial(l.Network, l.Address)
	}
	config, cErr := l.TLS.Config()
	if cErr != nil {
		return nil, cErr
	}
	return tls.Dial(l.Network, l.Address, config)
}
