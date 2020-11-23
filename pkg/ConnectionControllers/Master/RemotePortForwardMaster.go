package Master

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"strings"
	"time"
)

type RemotePortForward struct {
	Server                 net.Listener
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	TLSConfiguration       *tls.Config
	RemoteHost             string
	RemotePort             string
}

func (remotePortForward *RemotePortForward) Serve() error {
	var finalError error
	for {
		setTimeoutError := remotePortForward.MasterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		if setTimeoutError != nil {
			finalError = setTimeoutError
			break
		}
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(remotePortForward.MasterConnectionReader, 1)
		if connectionError != nil {
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				finalError = connectionError
				break
			}
		}
		if numberOfBytesReceived != 1 {
			_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &ConnectionControllers.UnknownOperation)
			continue
		}
		if buffer[0] != ConnectionControllers.NewConnection[0] {
			continue
		}
		targetConnection, connectionError := Sockets.Connect(&remotePortForward.RemoteHost, &remotePortForward.RemotePort)
		if connectionError != nil {
			_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &ConnectionControllers.FailToConnectToTarget)
			finalError = connectionError
			break
		}
		_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &ConnectionControllers.NewConnection)

		clientConnection, connectionError := remotePortForward.Server.Accept()
		clientConnection = Sockets.UpgradeServerToTLS(clientConnection, remotePortForward.TLSConfiguration)

		if connectionError != nil {
			// Do Something to notify the client that connection was not possible
			continue
		}
		if strings.Split(clientConnection.RemoteAddr().String(), ":")[0] == remotePortForward.RemoteHost {
			go ConnectionControllers.StartGeneralProxying(clientConnection, targetConnection)
		}
	}
	return finalError
}
