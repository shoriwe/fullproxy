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

type AuthenticationMethod func(username []byte, password []byte) bool

type ProxyProtocol interface {
	SetAuthenticationMethod(AuthenticationMethod) error
	Handle(net.Conn, *bufio.Reader, *bufio.Writer) error
}
