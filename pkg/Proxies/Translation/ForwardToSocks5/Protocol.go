package ForwardToSocks5

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"golang.org/x/net/proxy"
	"net"
)

type ForwardToSocks5 struct {
	TargetHost   string
	TargetPort   string
	Socks5Dialer proxy.Dialer
}

func (forwardToSocks5 *ForwardToSocks5) SetAuthenticationMethod(authenticationMethod ConnectionControllers.AuthenticationMethod) error {
	return nil
}

func (forwardToSocks5 *ForwardToSocks5) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	targetConnection, connectionError := forwardToSocks5.Socks5Dialer.Dial("tcp", forwardToSocks5.TargetHost+":"+forwardToSocks5.TargetPort)
	if connectionError == nil {
		targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		portProxy := Basic.PortProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		return portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
	_ = clientConnection.Close()
	return connectionError
}
