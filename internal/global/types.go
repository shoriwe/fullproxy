package global

import (
	"net"
)

type Protocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn) error
	SetOutboundFilter(IOFilter) error
	SetDial(dialFunc DialFunc)
	SetListen(listenFunc ListenFunc)
	SetListenAddress(address net.Addr)
}

type Pipe interface {
	SetLoggingMethod(LoggingMethod) error
	SetInboundFilter(IOFilter) error
	Serve() error
}
