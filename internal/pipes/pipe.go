package pipes

import (
	"net"
)

type (
	IOFilter      func(addr net.Addr) error
	LoggingMethod func(args ...interface{})
	Pipe          interface {
		LogData(a ...interface{})
		SetLoggingMethod(LoggingMethod)
		SetInboundFilter(IOFilter)
		SetOutboundFilter(IOFilter)
		FilterInbound(addr net.Addr) error
		FilterOutbound(addr net.Addr) error
		Serve() error
	}
)
