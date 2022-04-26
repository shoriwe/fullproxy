package http

import (
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
	"strconv"
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
	return &customListener{make(chan net.Conn, strconv.IntSize)}
}

type HTTP struct {
	AuthenticationMethod global.AuthenticationMethod
	proxyHttpServer      *goproxy.ProxyHttpServer
	LoggingMethod        global.LoggingMethod
	listener             *customListener
	OutboundFilter       global.IOFilter
}

func (protocol *HTTP) SetDial(dialFunc global.DialFunc) {
	protocol.proxyHttpServer.Tr.Dial = dialFunc
	protocol.proxyHttpServer.ConnectDial = dialFunc
}

func NewHTTP(
	authenticationMethod global.AuthenticationMethod,
	loggingMethod global.LoggingMethod,
	outboundFilter global.IOFilter,
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

				authenticationError := result.AuthenticationMethod(splitRawUsernamePassword[0], splitRawUsernamePassword[1])
				if authenticationError != nil {
					return req, onErrorResponse
				}
			}

			filterError := global.FilterOutbound(result.OutboundFilter, req.Host)
			if filterError != nil {
				return req, onErrorResponse
			}

			return req, nil
		},
	)
	return result
}

func (protocol *HTTP) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	protocol.LoggingMethod = loggingMethod
	return nil
}

func (protocol *HTTP) SetOutboundFilter(filter global.IOFilter) error {
	protocol.OutboundFilter = filter
	return nil
}

func (protocol *HTTP) SetAuthenticationMethod(authenticationMethod global.AuthenticationMethod) error {
	protocol.AuthenticationMethod = authenticationMethod
	return nil
}

func (protocol *HTTP) Handle(clientConnection net.Conn) error {
	protocol.listener.clientConnections <- clientConnection
	return nil
}
