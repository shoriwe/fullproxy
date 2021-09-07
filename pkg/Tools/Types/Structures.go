package Types

import (
	"net"
)

type ProxyProtocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn) error
	SetOutboundFilter(IOFilter) error
}

type Pipe interface {
	SetLoggingMethod(LoggingMethod) error
	SetInboundFilter(IOFilter) error
	Serve()
}
