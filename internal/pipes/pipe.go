package pipes

import (
	"net"
)

type (
	IOFilter      func(address string) error
	LoggingMethod func(args ...interface{})
	Pipe          interface {
		LogData(a ...interface{})
		SetLoggingMethod(LoggingMethod)
		SetInboundFilter(IOFilter)
		SetOutboundFilter(IOFilter)
		FilterInbound(address string) error
		FilterOutbound(address string) error
		Dial(networkType, address string) (net.Conn, error)
		Listen(networkType, address string) (net.Listener, error)
		Serve() error
	}
)
