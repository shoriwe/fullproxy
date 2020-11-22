package PortForward

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type RemoteForward struct {
	MasterHost string
	MasterPort string
	TLSConfiguration *tls.Config
}

func (remoteForward *RemoteForward)Handle(
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
			TargetConnection: targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		return portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	_ = clientConnection.Close()
	return connectionError
}

func (remoteForward *RemoteForward)SetAuthenticationMethod(_ ConnectionHandlers.AuthenticationFunction) error {
	return nil
}

func StartRemotePortForward(localAddress *string, localPort *string, masterAddress *string, masterPort *string) {
	if !(*localAddress == "" || *localPort == "" || *masterAddress == "" || *masterPort == "") {
		remoteForward := RemoteForward{
			MasterHost: *localAddress,
			MasterPort: *localPort,
		}
		ConnectionHandlers.RemotePortForwardSlave(masterAddress, masterPort, localAddress, localPort, &remoteForward)
	} else {
		fmt.Println("All flags need to be in use")
	}
}
