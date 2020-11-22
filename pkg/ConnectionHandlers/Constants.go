package ConnectionHandlers

import (
	"bufio"
	"crypto/tls"
	"net"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

type MasterFunction func(server net.Listener, masterConnection net.Conn, tlsConfiguration *tls.Config, args interface{})

type ProxyFunction func(conn net.Conn, connReader *bufio.Reader, connWriter *bufio.Writer, args ...interface{})

type AuthenticationFunction func(username []byte, password []byte) bool

type ProxyProtocol interface{
	SetAuthenticationMethod(AuthenticationFunction) error
	Handle(net.Conn, *bufio.Reader, *bufio.Writer) error
}