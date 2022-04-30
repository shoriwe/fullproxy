package pipes

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

type Bind struct {
	NetworkType                   string
	BindAddress                   string
	Server                        net.Listener
	Protocol                      servers.Protocol
	LoggingMethod                 LoggingMethod
	InboundFilter, OutboundFilter IOFilter
}

func (bind *Bind) LogData(a ...interface{}) {
	if bind.LoggingMethod != nil {
		bind.LoggingMethod(a...)
	}
}

func (bind *Bind) SetLoggingMethod(method LoggingMethod) {
	bind.LoggingMethod = method
}

func (bind *Bind) SetInboundFilter(filter IOFilter) {
	bind.InboundFilter = filter
}

func (bind *Bind) SetOutboundFilter(filter IOFilter) {
	bind.OutboundFilter = filter
}

func (bind *Bind) FilterInbound(addr net.Addr) error {
	if bind.InboundFilter != nil {
		return bind.InboundFilter(addr)
	}
	return nil
}

func (bind *Bind) FilterOutbound(addr net.Addr) error {
	if bind.OutboundFilter != nil {
		return bind.OutboundFilter(addr)
	}
	return nil
}

func (bind *Bind) serve(clientConnection net.Conn) {
	if filterError := bind.FilterInbound(clientConnection.RemoteAddr()); filterError != nil {
		bind.LogData(filterError)
		return
	}
	bind.LogData("Client connection received from: ", clientConnection.RemoteAddr().String())
	handleError := bind.Protocol.Handle(clientConnection)
	if handleError != nil {
		bind.LogData(handleError.Error())
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
	bind.LogData("Successfully listening at: " + bind.BindAddress)
	bind.Protocol.SetDial(net.Dial)     // TODO: Add filter inbound
	bind.Protocol.SetListen(net.Listen) // TODO: Add filter outbound
	bind.Protocol.SetListenAddress(listener.Addr())
	for {
		clientConnection, connectionError := bind.Server.Accept()
		if connectionError != nil {
			bind.LogData(connectionError)
			continue
		}
		go bind.serve(clientConnection)
	}
}

func NewBindPipe(
	networkType, bindAddress string,
	protocol servers.Protocol,
	method LoggingMethod,
	inboundFilter, outboundFilter IOFilter,
) Pipe {
	return &Bind{
		NetworkType:    networkType,
		BindAddress:    bindAddress,
		Protocol:       protocol,
		LoggingMethod:  method,
		InboundFilter:  inboundFilter,
		OutboundFilter: outboundFilter,
	}
}
