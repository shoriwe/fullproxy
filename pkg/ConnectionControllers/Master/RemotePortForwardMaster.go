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
	LoggingMethod          ConnectionControllers.LoggingMethod
	Tries                  int
	Timeout                time.Duration
}

func (remotePortForward *RemotePortForward) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	remotePortForward.LoggingMethod = loggingMethod
	return nil
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
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, connectionError)
			continue
		}
		if strings.Split(clientConnection.RemoteAddr().String(), ":")[0] == remotePortForward.RemoteHost {
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
			go ConnectionControllers.StartGeneralProxying(
				clientConnection, targetConnection,
				ConnectionControllers.GetTries(remotePortForward.Tries), ConnectionControllers.GetTimeout(remotePortForward.Timeout))
		} else {
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, "(Ignoring) Connection received from a non slave client: ", clientConnection.RemoteAddr().String())
		}
	}
	ConnectionControllers.LogData(remotePortForward.LoggingMethod, finalError)
	return finalError
}
