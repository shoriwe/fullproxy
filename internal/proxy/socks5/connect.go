package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/fullproxy/v3/internal/global"
	"github.com/shoriwe/fullproxy/v3/internal/pipes"
	"net"
)

func (socks5 *Socks5) Connect(context *Context) error {
	filterError := global.FilterOutbound(socks5.OutboundFilter, context.DSTAddress)
	if filterError != nil {
		_ = context.Reply(CommandReply{
			Version:    SocksV5,
			StatusCode: ConnectionNotAllowedByRuleSet,
			Address:    context.DSTRawAddress,
			Port:       context.DSTPort,
		})
		return filterError
	}

	// Try to connect to the target

	targetConnection, connectionError := socks5.Dial("tcp", context.DST)
	if connectionError != nil {
		_ = context.Reply(CommandReply{
			Version:    SocksV5,
			StatusCode: GeneralSocksServerFailure,
			Address:    context.DSTRawAddress,
			Port:       context.DSTPort,
		})
		return connectionError
	}

	// Respond to client

	targetConnectionAddress := targetConnection.LocalAddr().(*net.TCPAddr)
	var (
		bndAddress = targetConnectionAddress.IP
		bndPort    [2]byte
	)
	binary.BigEndian.PutUint16(bndPort[:], uint16(targetConnectionAddress.Port))

	if bndAddress.To4() != nil {
		bndAddress = bndAddress.To4()
	} else if bndAddress.To16() != nil {
		bndAddress = bndAddress.To16()
	}

	_ = context.Reply(CommandReply{
		Version:    SocksV5,
		StatusCode: ConnectionSucceed,
		Address:    bndAddress,
		Port:       context.DSTPort,
	})

	// Forward traffic
	return pipes.ForwardTraffic(context.ClientConnection, targetConnection)
}
