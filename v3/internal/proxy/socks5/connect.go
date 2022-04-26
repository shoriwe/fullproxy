package socks5

import (
	"encoding/binary"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"github.com/shoriwe/FullProxy/v3/internal/pipes"
	"net"
)

func (socks5 *Socks5) Connect(context *Context) error {
	filterError := global.FilterOutbound(socks5.OutboundFilter, context.DSTHost)
	if filterError != nil {
		_ = context.Reply(CommandReply{
			Version:     SocksV5,
			StatusCode:  ConnectionNotAllowedByRuleSet,
			AddressType: context.DSTAddressType,
			Address:     context.DSTRawAddress,
			Port:        context.DSTRawPort,
		})
		return filterError
	}

	// Try to connect to the target

	targetConnection, connectionError := socks5.Dial("tcp", context.DSTAddress)
	if connectionError != nil {
		_ = context.Reply(CommandReply{
			Version:     SocksV5,
			StatusCode:  GeneralSocksServerFailure,
			AddressType: context.DSTAddressType,
			Address:     context.DSTRawAddress,
			Port:        context.DSTRawPort,
		})
		return connectionError
	}

	// Respond to client

	targetConnectionAddress := targetConnection.LocalAddr().(*net.TCPAddr)
	var (
		bndType    byte
		bndAddress = targetConnectionAddress.IP
		bndPort    [2]byte
	)
	binary.BigEndian.PutUint16(bndPort[:], uint16(targetConnectionAddress.Port))

	if bndAddress.To4() != nil {
		bndAddress = bndAddress.To4()
		bndType = IPv4
	} else if bndAddress.To16() != nil {
		bndAddress = bndAddress.To16()
		bndType = IPv6
	}

	_ = context.Reply(CommandReply{
		Version:     SocksV5,
		StatusCode:  ConnectionSucceed,
		AddressType: bndType,
		Address:     bndAddress,
		Port:        bndPort[:],
	})

	// Forward traffic
	return pipes.ForwardTraffic(context.ClientConnection, targetConnection)
}
