package socks5

import (
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"net"
	"time"
)

func (socks5 *Socks5) Bind(context *Context) error {
	listener, listenError := socks5.Listen(
		"tcp",
		net.JoinHostPort(
			socks5.ListenAddress.IP.String(),
			"0",
		),
	)
	if listenError != nil {
		_ = context.Reply(
			CommandReply{
				Version:    SocksV5,
				StatusCode: ConnectionNotAllowedByRuleSet,
				Address:    context.ClientConnection.LocalAddr().(*net.TCPAddr).IP,
				Port:       0,
			},
		)
		return listenError
	}
	defer listener.Close()
	replyError := context.Reply(
		CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionSucceed,
			Address:    listener.Addr().(*net.TCPAddr).IP,
			Port:       listener.Addr().(*net.TCPAddr).Port,
		},
	)
	if replyError != nil {
		return replyError
	}
	switch listener.(type) {
	case *net.TCPListener:
		_ = listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Minute))
	case *listeners.TCPListener:
		_ = listener.(*listeners.TCPListener).Listener.SetDeadline(time.Now().Add(time.Minute))
	default:
		return errors.New("unsupported listener")
	}
	targetConnection, acceptError := listener.Accept()
	if acceptError != nil {
		_ = context.Reply(CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionNotAllowedByRuleSet,
			Address:    listener.Addr().(*net.TCPAddr).IP,
			Port:       listener.Addr().(*net.TCPAddr).Port,
		})
		return acceptError
	}
	defer targetConnection.Close()
	_ = listener.Close()
	if targetConnection.RemoteAddr().String() != context.DST {
		_ = context.Reply(CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionNotAllowedByRuleSet,
			Address:    targetConnection.RemoteAddr().(*net.TCPAddr).IP,
			Port:       targetConnection.RemoteAddr().(*net.TCPAddr).Port,
		})
		return ConnectionToReservedPort
	}
	_ = context.Reply(
		CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionSucceed,
			Address:    targetConnection.RemoteAddr().(*net.TCPAddr).IP,
			Port:       targetConnection.RemoteAddr().(*net.TCPAddr).Port,
		},
	)
	return common.ForwardTraffic(
		context.ClientConnection,
		&common.Sniffer{
			WriteSniffer: socks5.OutgoingSniffer,
			ReadSniffer:  socks5.IncomingSniffer,
			Connection:   targetConnection,
		},
	)
}
