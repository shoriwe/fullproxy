package listeners

import (
	"net"
)

type (
	LogFunc func(args ...interface{})
	Filters interface {
		Inbound(address string) error
		Outbound(address string) error
		Listen(address string) error
		Accept(address string) error
	}
	Listener interface {
		net.Listener
		Init() error
		Filter() Filters
		SetFilters(filters Filters)
		Dial(networkType, address string) (net.Conn, error)
		Listen(networkType, address string) (net.Listener, error)
	}
)
