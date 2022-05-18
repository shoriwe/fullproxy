package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"gopkg.in/elazarl/goproxy.v1"
	"io"
	"net"
	"net/http"
	"strings"
)

type HTTP struct {
	*goproxy.ProxyHttpServer
	AuthenticationMethod             servers.AuthenticationMethod
	IncomingSniffer, OutgoingSniffer io.Writer
}

func (H *HTTP) SetSniffers(incoming, outgoing io.Writer) {
	H.IncomingSniffer = incoming
	H.OutgoingSniffer = outgoing
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

func (H *HTTP) DoFunc(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	onErrorResponse := goproxy.NewResponse(req,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Don't waste your time!")
	newRequest := req.Clone(req.Context())
	newRequest.Body = &common.RequestSniffer{
		HeaderDone: false,
		Writer:     H.OutgoingSniffer,
		Request:    req,
	}
	req = newRequest
	if H.AuthenticationMethod != nil {
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

		authenticationError := H.AuthenticationMethod(splitRawUsernamePassword[0], splitRawUsernamePassword[1])
		if authenticationError != nil {
			return req, onErrorResponse
		}
	}

	return req, nil
}

func (H *HTTP) ResponseHandler(resp *http.Response, _ *goproxy.ProxyCtx) *http.Response {
	return &http.Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		ProtoMajor: resp.ProtoMajor,
		ProtoMinor: resp.ProtoMinor,
		Header:     resp.Header.Clone(),
		Body: &common.ResponseSniffer{
			HeaderDone: false,
			Writer:     H.IncomingSniffer,
			Response:   resp,
		},
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          resp.Trailer.Clone(),
		Request:          resp.Request,
		TLS:              resp.TLS,
	}
}

func NewHTTP(
	authenticationMethod servers.AuthenticationMethod,
) servers.HTTPHandler {
	protocol := &HTTP{
		ProxyHttpServer:      goproxy.NewProxyHttpServer(),
		AuthenticationMethod: authenticationMethod,
	}
	protocol.ProxyHttpServer.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	protocol.ProxyHttpServer.OnRequest().DoFunc(protocol.DoFunc)
	protocol.ProxyHttpServer.OnResponse().DoFunc(protocol.ResponseHandler)
	return protocol
}
