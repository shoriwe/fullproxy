package ForwardToSocks5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

type ForwardToSocks5 struct {
	TargetHost    string
	TargetPort    string
	Socks5Dialer  proxy.Dialer
	LoggingMethod ConnectionControllers.LoggingMethod
	Tries int
	Timeout time.Duration
}

func (forwardToSocks5 *ForwardToSocks5) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	forwardToSocks5.LoggingMethod = loggingMethod
	return nil
}
func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(authenticationMethod ConnectionControllers.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5)SetTries(tries int) error {
	forwardToSocks5.Tries = tries
	return nil
}

func (forwardToSocks5 *ForwardToSocks5)SetTimeout(timeout time.Duration) error {
	forwardToSocks5.Timeout = timeout
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetHost+":"+forwardToSocks5.TargetPort)
	if connectionError != nil {
		ConnectionControllers.LogData(forwardToSocks5.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	rawProxy := RawProxy.RawProxy{
		TargetConnection:       targetConnection,
		TargetConnectionReader: targetConnectionReader,
		TargetConnectionWriter: targetConnectionWriter,
		Tries: forwardToSocks5.Tries,
		Timeout: forwardToSocks5.Timeout,
	}
	return rawProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
}
