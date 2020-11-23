package ConnectionControllers

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

func LogData(loggingMethod LoggingMethod, arguments...interface{}) {
	if loggingMethod != nil {
		loggingMethod(arguments...)
	}
}

type AuthenticationMethod func(username []byte, password []byte) bool

type LoggingMethod func(args...interface{})

type ConnectionController interface {
	SetLoggingMethod(LoggingMethod)
	Serve() error
}

type ProxyProtocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn, *bufio.Reader, *bufio.Writer) error
}

func StartGeneralProxying(clientConnection net.Conn, targetConnection net.Conn) {
	clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	if targetConnectionReader != nil && targetConnectionWriter != nil {
		portProxy := PortProxy.PortProxy{
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
