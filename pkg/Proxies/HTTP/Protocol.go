package HTTP

import (
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
	"strings"
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

func NewHTTP(
	authenticationMethod Types.AuthenticationMethod,
	loggingMethod Types.LoggingMethod,
	outboundFilter Types.IOFilter,
) *HTTP {
	proxyHttpServer := goproxy.NewProxyHttpServer()
	listener := newCustomListener()
	go http.Serve(listener, proxyHttpServer)
	result := &HTTP{
		AuthenticationMethod: authenticationMethod,
		proxyHttpServer:      proxyHttpServer,
		LoggingMethod:        loggingMethod,
		listener:             listener,
		OutboundFilter:       outboundFilter,
	}
	proxyHttpServer.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxyHttpServer.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			onErrorResponse := goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusForbidden,
				"Don't waste your time!")
			if result.AuthenticationMethod != nil {

				entry := req.Header.Get("Proxy-Authorization")
				splitEntry := strings.Split(entry, " ")
				if len(splitEntry) != 2 {
					return req, onErrorResponse
				}

				if splitEntry[0] != "Basic" {
					return req, onErrorResponse
				}

				rawUsernamePassword, decodingError := base64.StdEncoding.DecodeString(splitEntry[1])
				if decodingError != nil {
					return req, onErrorResponse
				}

				splitRawUsernamePassword := bytes.Split(rawUsernamePassword, []byte{':'})
				if len(splitRawUsernamePassword) != 2 {
					return req, onErrorResponse
				}

				authentication, authenticationError := result.AuthenticationMethod(splitRawUsernamePassword[0], splitRawUsernamePassword[1])
				if !authentication || authenticationError != nil {
					return req, onErrorResponse
				}
			}

			if !Tools.FilterOutbound(result.OutboundFilter, req.Host) {
				return req, onErrorResponse
			}

			return req, nil
		},
	)
	return result
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
	protocol.AuthenticationMethod = authenticationMethod
	return nil
}

func (protocol *HTTP) Handle(clientConnection net.Conn) error {
	protocol.listener.clientConnections <- clientConnection
	return nil
}
