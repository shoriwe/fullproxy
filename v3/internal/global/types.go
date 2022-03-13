package global

import (
	"net"
)

type ProxyProtocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn) error
	SetOutboundFilter(IOFilter) error
	SetDial(dialFunc DialFunc)
}

type Pipe interface {
	SetLoggingMethod(LoggingMethod) error
	SetInboundFilter(IOFilter) error
	Serve() error
}
