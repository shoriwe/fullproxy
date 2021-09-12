package HTTP

import (
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
)

type customListener struct {
	clientConnections chan net.Conn
}

func (c customListener) Accept() (net.Conn, error) {
	return <-c.clientConnections, nil
}

func (c customListener) Close() error {
	close(c.clientConnections)
	return nil
}

func (c customListener) Addr() net.Addr {
	return nil
}

func newCustomListener() *customListener {
	return &customListener{make(chan net.Conn)}
}

type HTTP struct {
	AuthenticationMethod Types.AuthenticationMethod
	proxyHttpServer      *goproxy.ProxyHttpServer
	LoggingMethod        Types.LoggingMethod
	listener             *customListener
	OutboundFilter       Types.IOFilter
}

func (protocol *HTTP) SetDial(dialFunc Types.DialFunc) {
	protocol.proxyHttpServer.Tr.Dial = dialFunc
	protocol.proxyHttpServer.ConnectDial = dialFunc
}

func NewHTTP(authenticationMethod Types.AuthenticationMethod, loggingMethod Types.LoggingMethod, outboundFilter Types.IOFilter) *HTTP {
	proxyHttpServer := goproxy.NewProxyHttpServer()
	listener := newCustomListener()
	go http.Serve(listener, proxyHttpServer)
	return &HTTP{
		AuthenticationMethod: authenticationMethod,
		proxyHttpServer:      proxyHttpServer,
		LoggingMethod:        loggingMethod,
		listener:             listener,
		OutboundFilter:       outboundFilter,
	}
}

func (protocol *HTTP) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	protocol.LoggingMethod = loggingMethod
	return nil
}

func (protocol *HTTP) SetOutboundFilter(filter Types.IOFilter) error {
	protocol.OutboundFilter = filter
	return nil
}

func (protocol *HTTP) SetAuthenticationMethod(authenticationMethod Types.AuthenticationMethod) error {
	return nil
}

func (protocol *HTTP) Handle(clientConnection net.Conn) error {
	protocol.listener.clientConnections <- clientConnection
	return nil
}
