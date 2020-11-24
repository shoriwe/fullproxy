package Master

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
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
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
}

func (remotePortForward *RemotePortForward) SetInboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support InboundFilters")
}

func (remotePortForward *RemotePortForward) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of PIPE doesn't support OutboundFilters")
}

func (remotePortForward *RemotePortForward) SetTries(tries int) error {
	remotePortForward.Tries = tries
	return nil
}

func (remotePortForward *RemotePortForward) SetTimeout(timeout time.Duration) error {
	remotePortForward.Timeout = timeout
	return nil
}

func (remotePortForward *RemotePortForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
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
			_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &Pipes.UnknownOperation)
			continue
		}
		if buffer[0] != Pipes.NewConnection[0] {
			continue
		}
		targetConnection, connectionError := Sockets.Connect(&remotePortForward.RemoteHost, &remotePortForward.RemotePort)
		if connectionError != nil {
			_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &Pipes.FailToConnectToTarget)
			finalError = connectionError
			break
		}
		_, _ = Sockets.Send(remotePortForward.MasterConnectionWriter, &Pipes.NewConnection)

		clientConnection, connectionError := remotePortForward.Server.Accept()
		if connectionError != nil {
			_ = targetConnection.Close()
			// Do Something to notify the client that connection was not possible
			Templates.LogData(remotePortForward.LoggingMethod, connectionError)
			continue
		}
		clientConnection = Sockets.UpgradeServerToTLS(clientConnection, remotePortForward.TLSConfiguration)
		if strings.Split(clientConnection.RemoteAddr().String(), ":")[0] != remotePortForward.RemoteHost {
			Templates.LogData(remotePortForward.LoggingMethod, "(Ignoring) Connection received from a non slave client: ", clientConnection.RemoteAddr().String())
			_ = clientConnection.Close()
			_ = targetConnection.Close()
			continue
		}
		Templates.LogData(remotePortForward.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		go Pipes.StartGeneralProxying(
			clientConnection, targetConnection,
			Templates.GetTries(remotePortForward.Tries), Templates.GetTimeout(remotePortForward.Timeout))
	}
	Templates.LogData(remotePortForward.LoggingMethod, finalError)
	return finalError
}
