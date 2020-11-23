package HTTP

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

type CustomResponseWriter struct {
	httptest.ResponseRecorder
	http.Hijacker
	ClientConnection net.Conn
	ClientReadWriter *bufio.ReadWriter
}

func (customResponseWriter *CustomResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return customResponseWriter.ClientConnection, customResponseWriter.ClientReadWriter, nil
}

func CreateCustomResponseWriter(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) *CustomResponseWriter {
	slaveResponseWriter := new(CustomResponseWriter)
	slaveResponseWriter.Body = new(bytes.Buffer)
	slaveResponseWriter.Code = 200
	slaveResponseWriter.ClientConnection = clientConnection
	slaveResponseWriter.ClientReadWriter = bufio.NewReadWriter(
		clientConnectionReader,
		clientConnectionWriter)
	return slaveResponseWriter
}

type HTTP struct {
	AuthenticationMethod ConnectionControllers.AuthenticationMethod
	ProxyController      *goproxy.ProxyHttpServer
}

func (httpProtocol *HTTP) SetAuthenticationMethod(authenticationMethod ConnectionControllers.AuthenticationMethod) error {
	httpProtocol.ProxyController.OnRequest().DoFunc(func(
		request *http.Request,
		proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		authentication := strings.Split(request.Header.Get("Proxy-Authorization"), " ")
		if len(authentication) == 2 {
			if authentication[0] == "PortProxy" {
				rawCredentials, decodingError := base64.StdEncoding.DecodeString(authentication[1])
				if decodingError == nil {
					credentials := bytes.Split(rawCredentials, []byte(":"))
					if len(credentials) == 2 {
						if authenticationMethod(credentials[0], credentials[1]) {
							return request, nil
						}
					}
				}
			}
		}
		log.Print("Login failed with invalid credentials from: ", request.RemoteAddr)
		return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
	})
	return nil
}

func (httpProtocol *HTTP) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {

	request, parsingError := http.ReadRequest(clientConnectionReader)
	if parsingError == nil {
		request.RemoteAddr = clientConnection.RemoteAddr().String()
		responseWriter := CreateCustomResponseWriter(clientConnection, clientConnectionReader, clientConnectionWriter)
		httpProtocol.ProxyController.ServeHTTP(responseWriter, request)
		if responseWriter.Result().ContentLength > 0 {
			_ = responseWriter.Result().Write(clientConnectionWriter)
			_ = clientConnectionWriter.Flush()
		}
	}
	_ = clientConnection.Close()
	return parsingError
}
