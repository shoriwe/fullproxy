package PortForward

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type RemoteForward struct {
	MasterHost       string
	MasterPort       string
	TLSConfiguration *tls.Config
}

func (remoteForward *RemoteForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := Sockets.TLSConnect(
		&remoteForward.MasterHost,
		&remoteForward.MasterPort,
		(*remoteForward).TLSConfiguration)
	if connectionError == nil {
		targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		portProxy := Basic.PortProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		return portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	_ = clientConnection.Close()
	return connectionError
}

func (remoteForward *RemoteForward) SetAuthenticationMethod(_ ConnectionHandlers.AuthenticationMethod) error {
	return nil
}
