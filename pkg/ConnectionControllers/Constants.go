package ConnectionControllers

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
	"time"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

func LogData(loggingMethod LoggingMethod, arguments ...interface{}) {
	if loggingMethod != nil {
		loggingMethod(arguments...)
	}
}

func GetTries(tries int) int {
	if tries != 0 {
		return tries
	}
	return 5
}

func GetTimeout(timeout time.Duration) time.Duration {
	if timeout != 0 {
		return timeout
	}
	return 10 * time.Second
}

type AuthenticationMethod func(username []byte, password []byte) bool

type LoggingMethod func(args ...interface{})

type ConnectionController interface {
	SetLoggingMethod(LoggingMethod)
	Serve() error
	SetTries(int) error
	SetTimeout(time.Duration) error
}

type ProxyProtocol interface {
	SetLoggingMethod(LoggingMethod) error
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn, *bufio.Reader, *bufio.Writer) error
	SetTries(int) error
	SetTimeout(time.Duration) error
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
