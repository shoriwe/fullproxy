package RawProxy

import (
	"bufio"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

type RawProxy struct {
	TargetConnection       net.Conn
	TargetConnectionReader *bufio.Reader
	TargetConnectionWriter *bufio.Writer
	ConnectionAlive        bool
	LoggingMethod          Types.LoggingMethod
	Tries                  int
	Timeout                time.Duration
}

func (rawProxy *RawProxy) HandleReadWrite(
	sourceConnection net.Conn,
	sourceReader *bufio.Reader,
	destinationWriter *bufio.Writer) error {

	var proxyingError error
	tries := 0
	for ; tries < Templates.GetTries(rawProxy.Tries) && rawProxy.ConnectionAlive; tries++ {
		_ = sourceConnection.SetReadDeadline(time.Now().Add(Templates.GetTimeout(rawProxy.Timeout)))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(sourceReader, 1024)
		if connectionError != nil {
			// If the error is not "Timeout"
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				proxyingError = parsedConnectionError
				break
			}
		}
		tries = 0
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
	_ = sourceConnection.Close()
	if !rawProxy.ConnectionAlive {
		proxyingError = errors.New("connection died")
	} else {
		rawProxy.ConnectionAlive = false
	}
	if tries >= 5 {
		proxyingError = errors.New("max retries exceeded")
	}
	if proxyingError != nil {
		Templates.LogData(rawProxy.LoggingMethod, proxyingError)
	}
	return proxyingError
}

func (rawProxy *RawProxy) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}

func (rawProxy *RawProxy) SetInboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support InboundFilters")
}

func (rawProxy *RawProxy) SetOutboundFilter(_ Types.IOFilter) error {
	return errors.New("This kind of proxy doesn't support OutboundFilters")
}

func (rawProxy *RawProxy) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	rawProxy.LoggingMethod = loggingMethod
	return nil
}

func (rawProxy *RawProxy) SetTries(tries int) error {
	rawProxy.Tries = tries
	return nil
}

func (rawProxy *RawProxy) SetTimeout(timeout time.Duration) error {
	rawProxy.Timeout = timeout
	return nil
}

func (rawProxy *RawProxy) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	rawProxy.ConnectionAlive = true
	go rawProxy.HandleReadWrite(
		clientConnection,
		clientConnectionReader,
		rawProxy.TargetConnectionWriter)
	return rawProxy.HandleReadWrite(
		rawProxy.TargetConnection,
		rawProxy.TargetConnectionReader,
		clientConnectionWriter)
}
