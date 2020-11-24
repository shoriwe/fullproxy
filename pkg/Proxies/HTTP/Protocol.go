package HTTP

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Proxies"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"gopkg.in/elazarl/goproxy.v1"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	AuthenticationMethod Types.AuthenticationMethod
	ProxyController      *goproxy.ProxyHttpServer
	LoggingMethod        Types.LoggingMethod
	InboundFilter        Types.IOFilter
	OutboundFilter       Types.IOFilter
}

func (httpProtocol *HTTP) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	httpProtocol.LoggingMethod = loggingMethod
	return nil
}

func (httpProtocol *HTTP) SetInboundFilter(filter Types.IOFilter) error {
	if httpProtocol.ProxyController == nil {
		return errors.New("No HTTP proxy controller was set")
	}
	httpProtocol.InboundFilter = filter
	httpProtocol.ProxyController.OnRequest().DoFunc(func(
		request *http.Request,
		proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		clientRemoteAddressSplit := strings.Split(request.RemoteAddr, ":")
		clientRemoteAddress := new(net.TCPAddr)
		clientRemoteAddress.IP = net.IP(clientRemoteAddressSplit[0])
		portNumber, parsingError := strconv.Atoi(clientRemoteAddressSplit[1])
		if parsingError != nil {
			Templates.LogData(httpProtocol.LoggingMethod, parsingError)
			return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
		}
		clientRemoteAddress.Port = portNumber
		if !Templates.FilterInbound(httpProtocol.InboundFilter, clientRemoteAddress) {
			Templates.LogData(httpProtocol.LoggingMethod, "Connection denied to: ", request.RemoteAddr)
			return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
		}
		Templates.LogData(httpProtocol.LoggingMethod, "Connection Received from: ", request.RemoteAddr)
		return request, nil
	})
	return nil
}

func (httpProtocol *HTTP) SetOutboundFilter(filter Types.IOFilter) error {
	panic("Not implemented yet")
	/*
		if httpProtocol.ProxyController == nil {
			return errors.New("No HTTP proxy controller was set")
		}
	*/
	return nil
}

func (httpProtocol *HTTP) SetAuthenticationMethod(authenticationMethod Types.AuthenticationMethod) error {
	if httpProtocol.ProxyController == nil {
		return errors.New("No HTTP proxy controller was set")
	}
	httpProtocol.AuthenticationMethod = authenticationMethod
	httpProtocol.ProxyController.OnRequest().DoFunc(func(
		request *http.Request,
		proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		authentication := strings.Split(request.Header.Get("Proxy-Authorization"), " ")
		if len(authentication) == 2 {
			if authentication[0] == "Basic" {
				rawCredentials, decodingError := base64.StdEncoding.DecodeString(authentication[1])
				if decodingError == nil {
					credentials := bytes.Split(rawCredentials, []byte(":"))
					if len(credentials) == 2 {
						if Proxies.Authenticate(httpProtocol.AuthenticationMethod, credentials[0], credentials[1]) {
							Templates.LogData(httpProtocol.LoggingMethod, "Login successful with: ", request.RemoteAddr)
							return request, nil
						}
					}
				}
			}
		}
		Templates.LogData(httpProtocol.LoggingMethod, "Login failed with invalid credentials from: ", request.RemoteAddr)
		return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
	})
	return nil
}

func (httpProtocol *HTTP) SetTries(_ int) error {
	return errors.New("Custom tries is not supported by HTTP proxy")
}

func (httpProtocol *HTTP) SetTimeout(_ time.Duration) error {
	return errors.New("Custom timeout is not supported by HTTP proxy")
}

func (httpProtocol *HTTP) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {

	request, parsingError := http.ReadRequest(clientConnectionReader)
	if parsingError != nil {
		Templates.LogData(httpProtocol.LoggingMethod, parsingError)
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
