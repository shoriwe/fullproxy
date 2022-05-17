package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
	"strings"
)

type HTTP struct {
	*goproxy.ProxyHttpServer
}

func (H *HTTP) SetListenAddress(_ net.Addr) {
}

func (H *HTTP) SetListen(_ servers.ListenFunc) {
}

func (H *HTTP) Handle(_ net.Conn) error {
	panic("This should not be called")
}

func (H *HTTP) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (H *HTTP) SetDial(dialFunc servers.DialFunc) {
	H.ProxyHttpServer.Tr.DialContext = func(_ context.Context, network, addr string) (net.Conn, error) {
		return dialFunc(network, addr)
	}
}

func NewHTTP(
	authenticationMethod servers.AuthenticationMethod,
) servers.HTTPHandler {
	handler := goproxy.NewProxyHttpServer()
	handler.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	handler.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			onErrorResponse := goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusForbidden,
				"Don't waste your time!")
			if authenticationMethod != nil {

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

				authenticationError := authenticationMethod(splitRawUsernamePassword[0], splitRawUsernamePassword[1])
				if authenticationError != nil {
					return req, onErrorResponse
				}
			}

			return req, nil
		},
	)
	return &HTTP{
		ProxyHttpServer: handler,
	}
}
