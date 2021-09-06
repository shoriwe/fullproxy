package Types

import (
	"net"
	"time"
)

type ProxyProtocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn) error
	SetTries(int) error
	SetTimeout(time.Duration) error
	SetInboundFilter(IOFilter) error
	SetOutboundFilter(IOFilter) error
}

type Pipe interface {
	SetLoggingMethod(LoggingMethod) error
	SetInboundFilter(IOFilter) error
	SetOutboundFilter(IOFilter) error
	Serve() error
	SetTries(int) error
	SetTimeout(time.Duration) error
}
