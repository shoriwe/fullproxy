package PortForward

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type RemoteForward struct {
	MasterAddress    string
	MasterHost       string
	MasterPort       string
	TLSConfiguration *tls.Config
	LoggingMethod    Types.LoggingMethod
	Tries            int
	Timeout          time.Duration
	InboundFilter    Types.IOFilter
}

func (remoteForward *RemoteForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	remoteForward.LoggingMethod = loggingMethod
	return nil
}

func (remoteForward *RemoteForward) SetTries(tries int) error {
	remoteForward.Tries = tries
	return nil
}

func (remoteForward *RemoteForward) SetTimeout(timeout time.Duration) error {
	remoteForward.Timeout = timeout
	return nil
}

func (remoteForward *RemoteForward) SetInboundFilter(filter Types.IOFilter) error {
	remoteForward.InboundFilter = filter
	return nil
}

func (remoteForward *RemoteForward) Handle(clientConnection net.Conn) error {
	if !Tools.FilterInbound(remoteForward.InboundFilter, Tools.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
		Tools.LogData(remoteForward.LoggingMethod, errorMessage)
		_ = clientConnection.Close()
		return errors.New(errorMessage)
	}
	Tools.LogData(remoteForward.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	targetConnection, connectionError := Sockets.TLSConnect(
		remoteForward.MasterAddress,
		(*remoteForward).TLSConfiguration)
	if connectionError != nil {
		Tools.LogData(remoteForward.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}

func (remoteForward *RemoteForward) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}
