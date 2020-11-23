package HTTP

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
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
	LoggingMethod        ConnectionControllers.LoggingMethod
	// Tries int
	// Timeout time.Duration
}

func (httpProtocol *HTTP) SetLoggingMethod(loggingMethod ConnectionControllers.LoggingMethod) error {
	httpProtocol.LoggingMethod = loggingMethod
	return nil
}

func (httpProtocol *HTTP) SetAuthenticationMethod(authenticationMethod ConnectionControllers.AuthenticationMethod) error {
	if httpProtocol.ProxyController == nil {
		panic("No HTTP proxy controller was set")
	}
	httpProtocol.ProxyController.OnRequest().DoFunc(func(
		request *http.Request,
		proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		authentication := strings.Split(request.Header.Get("Proxy-Authorization"), " ")
		if len(authentication) == 2 {
			if authentication[0] == "RawProxy" {
				rawCredentials, decodingError := base64.StdEncoding.DecodeString(authentication[1])
				if decodingError == nil {
					credentials := bytes.Split(rawCredentials, []byte(":"))
					if len(credentials) == 2 {
						if Proxies.Authenticate(authenticationMethod, credentials[0], credentials[1]) {
							ConnectionControllers.LogData(httpProtocol.LoggingMethod, "Login successful with: ", request.RemoteAddr)
							return request, nil
						}
					}
				}
			}
		}
		ConnectionControllers.LogData(httpProtocol.LoggingMethod, "Login failed with invalid credentials from: ", request.RemoteAddr)
		return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
	})
	return nil
}

func (httpProtocol *HTTP) SetTries(_ int) error {
	// httpProtocol.Tries = tries
	return nil
}

func (httpProtocol *HTTP) SetTimeout(_ time.Duration) error {
	// httpProtocol.Timeout = timeout
	return nil
}

func (httpProtocol *HTTP) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {

	request, parsingError := http.ReadRequest(clientConnectionReader)
	if parsingError != nil {
		ConnectionControllers.LogData(httpProtocol.LoggingMethod, parsingError)
	} else {
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
