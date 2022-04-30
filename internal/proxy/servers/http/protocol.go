package http

import (
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
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
	AuthenticationMethod servers.AuthenticationMethod
	proxyHttpServer      *goproxy.ProxyHttpServer
	listener             *customListener
	ListenAddress        *net.TCPAddr
}

func (protocol *HTTP) SetListenAddress(address net.Addr) {
	protocol.ListenAddress = address.(*net.TCPAddr)
}

func (protocol *HTTP) SetListen(_ servers.ListenFunc) {
}

func (protocol *HTTP) SetDial(dialFunc servers.DialFunc) {
	protocol.proxyHttpServer.Tr.Dial = dialFunc
	protocol.proxyHttpServer.ConnectDial = dialFunc
}

func NewHTTP(
	authenticationMethod servers.AuthenticationMethod,
) servers.Protocol {
	proxyHttpServer := goproxy.NewProxyHttpServer()
	listener := newCustomListener()
	go http.Serve(listener, proxyHttpServer)
	result := &HTTP{
		AuthenticationMethod: authenticationMethod,
		proxyHttpServer:      proxyHttpServer,
		listener:             listener,
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

			return req, nil
		},
	)
	return result
}

func (protocol *HTTP) SetAuthenticationMethod(authenticationMethod servers.AuthenticationMethod) error {
	protocol.AuthenticationMethod = authenticationMethod
	return nil
}

func (protocol *HTTP) Handle(clientConnection net.Conn) error {
	protocol.listener.clientConnections <- clientConnection
	return nil
}
