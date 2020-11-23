package PortProxy

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

type PortProxy struct {
	TargetConnection       net.Conn
	TargetConnectionReader *bufio.Reader
	TargetConnectionWriter *bufio.Writer
	ConnectionAlive        bool
	Tries int
	Timeout time.Duration
}

func (portProxy *PortProxy)HandleReadWrite(
	sourceConnection net.Conn,
	sourceReader *bufio.Reader,
	destinationWriter *bufio.Writer,
	connectionAlive *bool) error {

	var proxyingError error
	tries := 0
	for ; tries < portProxy.Tries; tries++ {
		_ = sourceConnection.SetReadDeadline(time.Now().Add(portProxy.Timeout))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(sourceReader, 1048576)
		if connectionError != nil {
			// If the error is not "Timeout"
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				proxyingError = parsedConnectionError
				break
			} else {
				if !(*connectionAlive) {
					proxyingError = errors.New("connection died")
					break
				}
			}
		} else {
			tries = 0
		}
		if numberOfBytesReceived > 0 {
			realChunk := buffer[:numberOfBytesReceived]
			_, connectionError = Sockets.Send(destinationWriter, &realChunk)
			if connectionError != nil {
				proxyingError = connectionError
				break
			}
			realChunk = nil
		}
		buffer = nil
	}
	if tries >= 5 {
		proxyingError = errors.New("max retries exceeded")
	}
	_ = sourceConnection.Close()
	if *connectionAlive {
		*connectionAlive = false
	}
	return proxyingError
}

func (portProxy *PortProxy) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	portProxy.ConnectionAlive = true
	go portProxy.HandleReadWrite(
		clientConnection,
		clientConnectionReader,
		portProxy.TargetConnectionWriter,
		&portProxy.ConnectionAlive)
	return portProxy.HandleReadWrite(
		portProxy.TargetConnection,
		portProxy.TargetConnectionReader,
		clientConnectionWriter,
		&portProxy.ConnectionAlive)
}
