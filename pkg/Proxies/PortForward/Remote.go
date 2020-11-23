package PortForward

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type RemoteForward struct {
	MasterHost       string
	MasterPort       string
	TLSConfiguration *tls.Config
	LoggingMethod    ConnectionControllers.LoggingMethod
}

func (remoteForward *RemoteForward) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	remoteForward.LoggingMethod = loggingMethod
	return nil
}

func (remoteForward *RemoteForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := Sockets.TLSConnect(
		&remoteForward.MasterHost,
		&remoteForward.MasterPort,
		(*remoteForward).TLSConfiguration)
	if connectionError != nil {
		ConnectionControllers.LogData(remoteForward.LoggingMethod, connectionError)
	} else {
		targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		portProxy := PortProxy.PortProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		return portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	_ = clientConnection.Close()
	return connectionError
}

func (remoteForward *RemoteForward) SetAuthenticationMethod(_ ConnectionControllers.AuthenticationMethod) error {
	return nil
}
