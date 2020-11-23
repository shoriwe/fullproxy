package ConnectionControllers

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

type ConnectionController interface {
	Serve() error
}

type AuthenticationMethod func(username []byte, password []byte) bool

type ProxyProtocol interface {
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn, *bufio.Reader, *bufio.Writer) error
}

func StartGeneralProxying(clientConnection net.Conn, targetConnection net.Conn) {
	clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	if targetConnectionReader != nil && targetConnectionWriter != nil {
		portProxy := Basic.PortProxy{
			TargetConnection:       targetConnection,
			TargetConnectionReader: targetConnectionReader,
			TargetConnectionWriter: targetConnectionWriter,
		}
		portProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	} else {
		_ = clientConnection.Close()
		_ = targetConnection.Close()
	}
}
