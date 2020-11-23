package PortForward

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

type RemoteForward struct {
	MasterHost       string
	MasterPort       string
	TLSConfiguration *tls.Config
	LoggingMethod    ConnectionControllers.LoggingMethod
	Tries int
	Timeout time.Duration
}

func (remoteForward *RemoteForward) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
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
			Tries: ConnectionControllers.GetTries(remoteForward.Tries),
			Timeout: ConnectionControllers.GetTimeout(remoteForward.Timeout),
		}
		return portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	_ = clientConnection.Close()
	return connectionError
}

func (remoteForward *RemoteForward) SetAuthenticationMethod(_ ConnectionControllers.AuthenticationMethod) error {
	return nil
}
