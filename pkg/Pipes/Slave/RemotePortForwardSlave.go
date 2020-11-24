package Slave

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

type RemotePortForward struct {
	LocalServer            net.Listener
	MasterConnection       net.Conn
	MasterConnectionReader *bufio.Reader
	MasterConnectionWriter *bufio.Writer
	MasterHost             string
	MasterPort             string
	TLSConfiguration       *tls.Config
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
	InboundFilter          Types.IOFilter
}

func (remotePortForward *RemotePortForward) SetInboundFilter(filter Types.IOFilter) error {
	remotePortForward.InboundFilter = filter
	return nil
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
		clientConnection, connectionError := remotePortForward.LocalServer.Accept()
		if connectionError != nil {
			finalError = connectionError
			break
		}
		if !Templates.FilterInbound(remotePortForward.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
			errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
			_ = clientConnection.Close()
			Templates.LogData(remotePortForward.LoggingMethod, errorMessage)
			continue
		}
		Templates.LogData(remotePortForward.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		_, connectionError = Sockets.Send(remotePortForward.MasterConnectionWriter, &Pipes.NewConnection)
		if connectionError != nil {
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				finalError = connectionError
				break
			}
		}
		_ = remotePortForward.MasterConnection.SetReadDeadline(time.Now().Add(3 * time.Second))
		numberOfBytesReceived, response, connectionError := Sockets.Receive(remotePortForward.MasterConnectionReader, 1)
		if connectionError != nil {
			Templates.LogData(remotePortForward.LoggingMethod, connectionError)
			continue
		}
		if numberOfBytesReceived != 1 {
			continue
		}
		switch response[0] {
		case Pipes.NewConnection[0]:
			targetConnection, connectionError := Sockets.TLSConnect(
				&remotePortForward.MasterHost,
				&remotePortForward.MasterPort,
				remotePortForward.TLSConfiguration)
			if connectionError != nil {
				_ = clientConnection.Close()
				Templates.LogData(remotePortForward.LoggingMethod, "Connectivity issues with master server")
				return connectionError
			}
			go Pipes.StartGeneralProxying(
				clientConnection, targetConnection,
				Templates.GetTries(remotePortForward.Tries), Templates.GetTimeout(remotePortForward.Timeout))
		case Pipes.FailToConnectToTarget[0]:
			_ = clientConnection.Close()
			Templates.LogData(remotePortForward.LoggingMethod, "Something goes wrong when master connected to target")
			return errors.New("Something goes wrong when master connected to target")
		case Pipes.UnknownOperation[0]:
			_ = clientConnection.Close()
			Templates.LogData(remotePortForward.LoggingMethod, "The master did not understood the message")
			return errors.New("The master did not understood the message")
		}
	}
	return finalError
}
