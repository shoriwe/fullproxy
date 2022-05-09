package servers

import (
	"net"
	"net/http"
)

type (
	AuthenticationMethod func(username []byte, password []byte) error

	DialFunc func(network, address string) (net.Conn, error)

	DialUDPFunc func(network string, localAddress, remoteAddress *net.UDPAddr) (*net.UDPConn, error)

	ListenFunc func(network, address string) (net.Listener, error)

	Protocol interface {
		Handle(net.Conn) error
		SetAuthenticationMethod(AuthenticationMethod)
		SetDial(dialFunc DialFunc)
		SetListen(listenFunc ListenFunc)
		SetListenAddress(address net.Addr)
	}

	HTTPHandler interface {
		Protocol
		http.Handler
	}
)
