package proxies

import (
	"context"
	"net"

	"github.com/things-go/go-socks5"
)

type Socks5 struct {
	Listener    net.Listener
	Dial        func(network, addr string) (net.Conn, error)
	AuthMethods []socks5.Authenticator
}

func (s *Socks5) Serve() error {
	server := socks5.NewServer(
		socks5.WithDial(func(ctx context.Context, network, addr string) (net.Conn, error) {
			return s.Dial(network, addr)
		}),
		// socks5.WithConnectHandle(
		// 	func(ctx context.Context, writer io.Writer, request *socks5.Request) error {
		//
		// 	},
		// ),
		socks5.WithAuthMethods(s.AuthMethods),
		socks5.WithRule(&socks5.PermitCommand{EnableConnect: true}), // FIXME: Maybe in the future this could change
	)
	return server.Serve(s.Listener)
}
