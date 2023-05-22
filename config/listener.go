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
	ls, lErr := net.Listen(l.Network, l.Address)
	if lErr != nil {
		return nil, lErr
	}
	if l.TLS == nil {
		return ls, nil
	}
	config, cErr := l.TLS.Config()
	if cErr != nil {
		ls.Close()
		return nil, cErr
	}
	return tls.NewListener(ls, config), nil
}
