package socks5

import (
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
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
		return listenError
	}
	defer listener.Close()
	replyError := context.Reply(
		CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionSucceed,
			Address:    context.ClientConnection.LocalAddr().(*net.TCPAddr).IP,
			Port:       listener.Addr().(*net.TCPAddr).Port,
		},
	)
	if replyError != nil {
		return replyError
	}
	_ = listener.(*net.TCPListener).SetDeadline(time.Now().Add(30 * time.Second))
	targetConnection, acceptError := listener.Accept()
	if acceptError != nil {
		return acceptError
	}
	defer targetConnection.Close()
	_ = listener.Close()
	if !targetConnection.RemoteAddr().(*net.TCPAddr).IP.Equal(context.DSTRawAddress) {
		_ = context.Reply(CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionNotAllowedByRuleSet,
			Address:    targetConnection.RemoteAddr().(*net.TCPAddr).IP,
			Port:       targetConnection.RemoteAddr().(*net.TCPAddr).Port,
		})
		return ConnectionToReservedPort
	}
	replyError = context.Reply(
		CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionSucceed,
			Address:    targetConnection.RemoteAddr().(*net.TCPAddr).IP,
			Port:       targetConnection.RemoteAddr().(*net.TCPAddr).Port,
		},
	)
	return pipes.ForwardTraffic(context.ClientConnection, targetConnection)
}
