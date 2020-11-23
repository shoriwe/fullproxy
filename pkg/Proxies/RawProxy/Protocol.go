package RawProxy

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

type RawProxy struct {
	TargetConnection       net.Conn
	TargetConnectionReader *bufio.Reader
	TargetConnectionWriter *bufio.Writer
	ConnectionAlive        bool
	Tries int
	Timeout time.Duration
}

func (rawProxy *RawProxy)HandleReadWrite(
	sourceConnection net.Conn,
	sourceReader *bufio.Reader,
	destinationWriter *bufio.Writer,
	connectionAlive *bool) error {

	var proxyingError error
	tries := 0
	for ; tries < rawProxy.Tries; tries++ {
		_ = sourceConnection.SetReadDeadline(time.Now().Add(rawProxy.Timeout))
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

func (rawProxy *RawProxy) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	rawProxy.ConnectionAlive = true
	go rawProxy.HandleReadWrite(
		clientConnection,
		clientConnectionReader,
		rawProxy.TargetConnectionWriter,
		&rawProxy.ConnectionAlive)
	return rawProxy.HandleReadWrite(
		rawProxy.TargetConnection,
		rawProxy.TargetConnectionReader,
		clientConnectionWriter,
		&rawProxy.ConnectionAlive)
}
