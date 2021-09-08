package Pipes

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
)

type Bind struct {
	NetworkType   string
	BindAddress   string
	Server        net.Listener
	ProxyProtocol Types.ProxyProtocol
	LoggingMethod Types.LoggingMethod
	InboundFilter Types.IOFilter
}

func (bind *Bind) SetInboundFilter(filter Types.IOFilter) error {
	bind.InboundFilter = filter
	return nil
}

func (bind *Bind) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (bind *Bind) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	bind.LoggingMethod = loggingMethod
	return nil
}

func (bind *Bind) serve(clientConnection net.Conn) {
	if !Tools.FilterInbound(bind.InboundFilter, Tools.ParseIP(clientConnection.RemoteAddr().String()).String()) {
		_ = clientConnection.Close()
		Tools.LogData(bind.LoggingMethod, "Connection denied to: "+clientConnection.RemoteAddr().String())
		return
	}
	Tools.LogData(bind.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
	handleError := bind.ProxyProtocol.Handle(clientConnection)
	if handleError != nil {
		Tools.LogData(bind.LoggingMethod, handleError.Error)
	}
	return
}

func (bind *Bind) Serve() error {
	listener, listenError := net.Listen(bind.NetworkType, bind.BindAddress)
	if listenError != nil {
		return listenError
	}
	bind.Server = listener
	defer bind.Server.Close()
	Tools.LogData(bind.LoggingMethod, "Successfully listening at: "+bind.BindAddress)
	bind.ProxyProtocol.SetDial(net.Dial)
	for {
		clientConnection, connectionError := bind.Server.Accept()
		if connectionError != nil {
			Tools.LogData(bind.LoggingMethod, connectionError)
			continue
		}
		go bind.serve(clientConnection)
	}
}

func NewBindPipe(networkType, bindAddress string, protocol Types.ProxyProtocol, method Types.LoggingMethod, inboundFilter Types.IOFilter) *Bind {
	return &Bind{
		NetworkType:   networkType,
		BindAddress:   bindAddress,
		ProxyProtocol: protocol,
		LoggingMethod: method,
		InboundFilter: inboundFilter,
	}
}
