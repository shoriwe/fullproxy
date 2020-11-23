package Slave

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
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
		clientConnection, connectionError := remotePortForward.LocalServer.Accept()
		if connectionError != nil {
			finalError = connectionError
			break
		}
		ConnectionControllers.LogData(remotePortForward.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		_, connectionError = Sockets.Send(remotePortForward.MasterConnectionWriter, &ConnectionControllers.NewConnection)
		if connectionError != nil {
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				finalError = connectionError
				break
			}
		}
		_ = remotePortForward.MasterConnection.SetReadDeadline(time.Now().Add(3 * time.Second))
		numberOfBytesReceived, response, connectionError := Sockets.Receive(remotePortForward.MasterConnectionReader, 1)
		if connectionError != nil {
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, connectionError)
			continue
		}
		if numberOfBytesReceived != 1 {
			continue
		}
		switch response[0] {
		case ConnectionControllers.NewConnection[0]:
			targetConnection, connectionError := Sockets.TLSConnect(
				&remotePortForward.MasterHost,
				&remotePortForward.MasterPort,
				remotePortForward.TLSConfiguration)
			if connectionError == nil {
				go ConnectionControllers.StartGeneralProxying(
					clientConnection, targetConnection,
					ConnectionControllers.GetTries(remotePortForward.Tries), ConnectionControllers.GetTimeout(remotePortForward.Timeout))
			} else {
				_ = clientConnection.Close()
				ConnectionControllers.LogData(remotePortForward.LoggingMethod, "Connectivity issues with master server")
				return connectionError
			}
		case ConnectionControllers.FailToConnectToTarget[0]:
			_ = clientConnection.Close()
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, "Something goes wrong when master connected to target")
			return errors.New("Something goes wrong when master connected to target")
		case ConnectionControllers.UnknownOperation[0]:
			_ = clientConnection.Close()
			ConnectionControllers.LogData(remotePortForward.LoggingMethod, "The master did not understood the message")
			return errors.New("The master did not understood the message")
		}
	}
	return finalError
}
