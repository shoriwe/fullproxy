package ForwardToSocks5

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

type ForwardToSocks5 struct {
	TargetHost    string
	TargetPort    string
	Socks5Dialer  proxy.Dialer
	LoggingMethod Types.LoggingMethod
	Tries         int
	Timeout       time.Duration
	InboundFilter Types.IOFilter
}

func (forwardToSocks5 *ForwardToSocks5) SetInboundFilter(filter Types.IOFilter) error {
	forwardToSocks5.InboundFilter = filter
	return nil
}

func (ForwardToSocks5) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (forwardToSocks5 *ForwardToSocks5) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	forwardToSocks5.LoggingMethod = loggingMethod
	return nil
}
func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (forwardToSocks5 *ForwardToSocks5) SetTries(tries int) error {
	forwardToSocks5.Tries = tries
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) SetTimeout(timeout time.Duration) error {
	forwardToSocks5.Timeout = timeout
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(
	clientConnection net.Conn) error {
	if !Tools.FilterInbound(forwardToSocks5.InboundFilter, Tools.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
		Tools.LogData(forwardToSocks5.LoggingMethod, errorMessage)
		_ = clientConnection.Close()
		return errors.New(errorMessage)
	}
	Tools.LogData(forwardToSocks5.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetHost+":"+forwardToSocks5.TargetPort)
	if connectionError != nil {
		Tools.LogData(forwardToSocks5.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}
