package PortForward

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

type LocalForward struct {
	TargetHost    string
	TargetPort    string
	LoggingMethod Types.LoggingMethod
	Tries         int
	Timeout       time.Duration
	InboundFilter Types.IOFilter
}

func (localForward *LocalForward) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (localForward *LocalForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	localForward.LoggingMethod = loggingMethod
	return nil
}

func (localForward *LocalForward) SetTries(tries int) error {
	localForward.Tries = tries
	return nil
}

func (localForward *LocalForward) SetTimeout(timeout time.Duration) error {
	localForward.Timeout = timeout
	return nil
}

func (localForward *LocalForward) SetInboundFilter(filter Types.IOFilter) error {
	localForward.InboundFilter = filter
	return nil
}

func (localForward *LocalForward) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (localForward *LocalForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	if !Templates.FilterInbound(localForward.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
		Templates.LogData(localForward.LoggingMethod, errorMessage)
		_ = clientConnection.Close()
		return errors.New(errorMessage)
	}
	Templates.LogData(localForward.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	targetConnection, connectionError := Sockets.Connect(&localForward.TargetHost, &localForward.TargetPort)
	if connectionError != nil {
		Templates.LogData(localForward.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	targetReader, targetWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	rawProxy := RawProxy.RawProxy{
		TargetConnection:       targetConnection,
		TargetConnectionReader: targetReader,
		TargetConnectionWriter: targetWriter,
	}
	_ = rawProxy.SetTries(localForward.Tries)
	_ = rawProxy.SetTimeout(localForward.Timeout)
	_ = rawProxy.SetLoggingMethod(localForward.LoggingMethod)
	return rawProxy.Handle(
		clientConnection,
		clientConnectionReader, clientConnectionWriter,
	)
}
