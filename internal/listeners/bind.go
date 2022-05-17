package listeners

import (
	"crypto/tls"
	"net"
)

type BindListener struct {
	filters Filters
	net.Listener
}

func (bind *BindListener) Filter() Filters {
	return bind.filters
}

func (bind *BindListener) SetFilters(filters Filters) {
	bind.filters = filters
}

func (bind *BindListener) Init() error {
	return nil
}

func (bind *BindListener) Dial(networkType, address string) (net.Conn, error) {
	return net.Dial(networkType, address)
}

func (bind *BindListener) Listen(networkType, address string) (net.Listener, error) {
	return net.Listen(networkType, address)
}

func NewBindListener(
	networkType, address string,
	config *tls.Config,
) (Listener, error) {
	var (
		listener    net.Listener
		listenError error
	)
	if config == nil {
		listener, listenError = net.Listen(networkType, address)
	} else {
		listener, listenError = tls.Listen(networkType, address, config)
	}
	if listenError != nil {
		return nil, listenError
	}
	return &BindListener{
		Listener: listener,
	}, nil
}
