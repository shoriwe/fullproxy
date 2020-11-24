package ForwardToSocks5

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
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
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	if !Templates.FilterInbound(forwardToSocks5.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
		Templates.LogData(forwardToSocks5.LoggingMethod, errorMessage)
		_ = clientConnection.Close()
		return errors.New(errorMessage)
	}
	Templates.LogData(forwardToSocks5.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetHost+":"+forwardToSocks5.TargetPort)
	if connectionError != nil {
		Templates.LogData(forwardToSocks5.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	rawProxy := RawProxy.RawProxy{
		TargetConnection:       targetConnection,
		TargetConnectionReader: targetConnectionReader,
		TargetConnectionWriter: targetConnectionWriter,
	}
	_ = rawProxy.SetTries(forwardToSocks5.Tries)
	_ = rawProxy.SetTimeout(forwardToSocks5.Timeout)
	_ = rawProxy.SetLoggingMethod(forwardToSocks5.LoggingMethod)
	return rawProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
}
