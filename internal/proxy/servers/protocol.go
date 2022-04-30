package servers

import (
	"net"
)

type (
	AuthenticationMethod func(username []byte, password []byte) error

	DialFunc func(network, address string) (net.Conn, error)

	DialUDPFunc func(network string, localAddress, remoteAddress *net.UDPAddr) (*net.UDPConn, error)

	ListenFunc func(network, address string) (net.Listener, error)

	Protocol interface {
		SetAuthenticationMethod(AuthenticationMethod) error
		Handle(net.Conn) error
		SetDial(dialFunc DialFunc)
		SetListen(listenFunc ListenFunc)
		SetListenAddress(address net.Addr)
	}
)
