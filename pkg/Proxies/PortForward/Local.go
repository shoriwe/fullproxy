package PortForward

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type LocalForward struct {
	TargetHost    string
	TargetPort    string
	LoggingMethod ConnectionControllers.LoggingMethod
}

func (localForward *LocalForward) SetAuthenticationMethod(_ ConnectionControllers.AuthenticationMethod) error {
	return nil
}

func (localForward *LocalForward) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	localForward.LoggingMethod = loggingMethod
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
		portProxy := PortProxy.PortProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetReader,
			TargetConnectionWriter: targetWriter,
		}
		return portProxy.Handle(
			clientConnection,
			clientConnectionReader, clientConnectionWriter,
		)
	}
	_ = clientConnection.Close()
	return connectionError
}
