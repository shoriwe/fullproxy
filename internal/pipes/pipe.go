package pipes

import (
	"crypto/tls"
	"net"
)

type (
	IOFilter      func(address string) error
	LoggingMethod func(args ...interface{})
	Pipe          interface {
		LogData(a ...interface{})
		SetTLSCertificates(certificates []tls.Certificate)
		SetLoggingMethod(method LoggingMethod)
		SetInboundFilter(filter IOFilter)
		SetOutboundFilter(filter IOFilter)
		SetListenFilter(filter IOFilter)
		SetAcceptFilter(filter IOFilter)
		FilterInbound(address string) error
		FilterOutbound(address string) error
		FilterListen(address string) error
		Dial(networkType, address string) (net.Conn, error)
		Listen(networkType, address string) (net.Listener, error)
		Serve() error
	}
)
