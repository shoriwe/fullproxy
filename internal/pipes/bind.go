package pipes

import (
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/global"
	"net"
)

type Bind struct {
	NetworkType   string
	BindAddress   string
	Server        net.Listener
	Protocol      global.Protocol
	LoggingMethod global.LoggingMethod
	InboundFilter global.IOFilter
}

func (bind *Bind) SetInboundFilter(filter global.IOFilter) error {
	bind.InboundFilter = filter
	return nil
}

func (bind *Bind) SetOutboundFilter(_ global.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (bind *Bind) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	bind.LoggingMethod = loggingMethod
	return nil
}

func (bind *Bind) serve(clientConnection net.Conn) {
	filterError := global.FilterInbound(bind.InboundFilter, global.ParseIP(clientConnection.RemoteAddr().String()).String())
	if filterError != nil {
		_ = clientConnection.Close()
		global.LogData(bind.LoggingMethod, filterError)
		return
	}
	global.LogData(bind.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
	handleError := bind.Protocol.Handle(clientConnection)
	if handleError != nil {
		global.LogData(bind.LoggingMethod, handleError.Error())
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
	global.LogData(bind.LoggingMethod, "Successfully listening at: "+bind.BindAddress)
	bind.Protocol.SetDial(net.Dial)
	bind.Protocol.SetListen(net.Listen)
	bind.Protocol.SetListenAddress(listener.Addr())
	for {
		clientConnection, connectionError := bind.Server.Accept()
		if connectionError != nil {
			global.LogData(bind.LoggingMethod, connectionError)
			continue
		}
		go bind.serve(clientConnection)
	}
}

func NewBindPipe(networkType, bindAddress string, protocol global.Protocol, method global.LoggingMethod, inboundFilter global.IOFilter) *Bind {
	return &Bind{
		NetworkType:   networkType,
		BindAddress:   bindAddress,
		Protocol:      protocol,
		LoggingMethod: method,
		InboundFilter: inboundFilter,
	}
}
