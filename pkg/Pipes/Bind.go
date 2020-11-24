package Pipes

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

type Bind struct {
	Server        net.Listener
	ProxyProtocol Types.ProxyProtocol
	LoggingMethod Types.LoggingMethod
	InboundFilter Types.InboundFilter
}

func (bind *Bind) SetInboundFilter(filter Types.InboundFilter) error {
	bind.InboundFilter = filter
	return nil
}

func (bind *Bind) SetOutboundFilter(_ Types.OutboundFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (bind *Bind) SetTries(tries int) error {
	return bind.ProxyProtocol.SetTries(tries)
}

func (bind *Bind) SetTimeout(timeout time.Duration) error {
	return bind.ProxyProtocol.SetTimeout(timeout)
}

func (bind *Bind) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	bind.LoggingMethod = loggingMethod
	return nil
}

func (bind *Bind) Serve() error {
	for {
		clientConnection, connectionError := bind.Server.Accept()
		if connectionError != nil {
			Templates.LogData(bind.LoggingMethod, connectionError)
			return connectionError
		}
		if !Templates.FilterInbound(bind.InboundFilter, clientConnection.RemoteAddr()) {
			Templates.LogData(bind.LoggingMethod, "Unwanted connection received from "+clientConnection.RemoteAddr().String())
			continue
		}
		Templates.LogData(bind.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go bind.ProxyProtocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
}
