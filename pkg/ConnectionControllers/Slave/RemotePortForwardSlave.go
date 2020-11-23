package Slave

import (
	"bufio"
	"crypto/tls"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
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
}

func (remotePortForward *RemotePortForward) Serve() error {
	var finalError error
	for {
		clientConnection, connectionError := remotePortForward.LocalServer.Accept()
		if connectionError != nil {
			finalError = connectionError
			break
		}
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
				go ConnectionControllers.StartGeneralProxying(clientConnection, targetConnection)
			} else {
				_ = clientConnection.Close()
				log.Fatal("Connectivity issues with master server")
			}
		case ConnectionControllers.FailToConnectToTarget[0]:
			_ = clientConnection.Close()
			log.Print("Something goes wrong when master connected to target")
			break
		case ConnectionControllers.UnknownOperation[0]:
			_ = clientConnection.Close()
			log.Print("The master did not understood the message")
			break
		}
	}
	return finalError
}
