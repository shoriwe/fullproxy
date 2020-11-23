package PortProxy

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

func HandleReadWrite(
	sourceConnection net.Conn,
	sourceReader *bufio.Reader,
	destinationWriter *bufio.Writer,
	connectionAlive *bool) error {

	var proxyingError error
	tries := 0
	for ; tries < 5; tries++ {
		_ = sourceConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
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

type PortProxy struct {
	TargetConnection       net.Conn
	TargetConnectionReader *bufio.Reader
	TargetConnectionWriter *bufio.Writer
	ConnectionAlive        bool
}

func (portProxy *PortProxy) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	portProxy.ConnectionAlive = true
	go HandleReadWrite(
		clientConnection,
		clientConnectionReader,
		portProxy.TargetConnectionWriter,
		&portProxy.ConnectionAlive)
	return HandleReadWrite(
		portProxy.TargetConnection,
		portProxy.TargetConnectionReader,
		clientConnectionWriter,
		&portProxy.ConnectionAlive)
}
