package pipes

import (
	"crypto/tls"
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
	TLSCertificates               []tls.Certificate
}

func (bind *Bind) Dial(networkType, address string) (net.Conn, error) {
	if filterError := bind.FilterInbound(address); filterError != nil {
		return nil, filterError
	}
	return net.Dial(networkType, address)
}

func (bind *Bind) Listen(networkType, address string) (net.Listener, error) {
	// TODO: Listen filter
	// TODO: Return Filterable listener
	return net.Listen(networkType, address)
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

func (bind *Bind) FilterInbound(address string) error {
	if bind.InboundFilter != nil {
		return bind.InboundFilter(address)
	}
	return nil
}

func (bind *Bind) FilterOutbound(address string) error {
	if bind.OutboundFilter != nil {
		return bind.OutboundFilter(address)
	}
	return nil
}

func (bind *Bind) serve(clientConnection net.Conn) {
	if filterError := bind.FilterInbound(clientConnection.RemoteAddr().String()); filterError != nil {
		bind.LogData(filterError)
		return
	}
	bind.LogData("Client connection received from: ", clientConnection.RemoteAddr().String())
	handleError := bind.Protocol.Handle(clientConnection)
	if handleError != nil {
		bind.LogData(handleError.Error())
	}
}

func (bind *Bind) Serve() error {
	var (
		listener    net.Listener
		listenError error
	)
	if bind.TLSCertificates == nil {
		listener, listenError = net.Listen(bind.NetworkType, bind.BindAddress)
	} else {
		listener, listenError = tls.Listen(bind.NetworkType, bind.BindAddress, &tls.Config{Certificates: bind.TLSCertificates})
	}
	if listenError != nil {
		return listenError
	}
	bind.Server = listener
	defer bind.Server.Close()
	bind.LogData("Successfully listening at: " + bind.BindAddress)
	bind.Protocol.SetDial(bind.Dial)
	bind.Protocol.SetListen(bind.Listen)
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
	certificates []tls.Certificate,
) Pipe {
	return &Bind{
		NetworkType:     networkType,
		BindAddress:     bindAddress,
		Protocol:        protocol,
		LoggingMethod:   method,
		InboundFilter:   inboundFilter,
		OutboundFilter:  outboundFilter,
		TLSCertificates: certificates,
	}
}
