package PortForward

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type LocalForward struct {
	TargetHost string
	TargetPort string
}

func (localForward *LocalForward) SetAuthenticationMethod(_ ConnectionControllers.AuthenticationMethod) error {
	return nil
}

func (localForward *LocalForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := Sockets.Connect(&localForward.TargetHost, &localForward.TargetPort)
	if connectionError == nil {
		targetReader, targetWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		portProxy := Basic.PortProxy{
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
