package PortForward

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

type LocalForward struct {
	TargetHost    string
	TargetPort    string
	LoggingMethod ConnectionControllers.LoggingMethod
	Tries         int
	Timeout       time.Duration
}

func (localForward *LocalForward) SetAuthenticationMethod(_ ConnectionControllers.AuthenticationMethod) error {
	return nil
}

func (localForward *LocalForward) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
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

func (localForward *LocalForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := Sockets.Connect(&localForward.TargetHost, &localForward.TargetPort)
	if connectionError != nil {
		ConnectionControllers.LogData(localForward.LoggingMethod, connectionError)
	} else {
		targetReader, targetWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		rawProxy := RawProxy.RawProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetReader,
			TargetConnectionWriter: targetWriter,
			Tries:                  ConnectionControllers.GetTries(localForward.Tries),
			Timeout:                ConnectionControllers.GetTimeout(localForward.Timeout),
		}
		return rawProxy.Handle(
			clientConnection,
			clientConnectionReader, clientConnectionWriter,
		)
	}
	_ = clientConnection.Close()
	return connectionError
}
