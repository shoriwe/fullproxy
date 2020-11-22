package HTTP

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/shoriwe/FullProxy/pkg/Hashing"
	"github.com/shoriwe/FullProxy/pkg/MasterSlave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

type SlaveResponseWriter struct {
	httptest.ResponseRecorder
	http.Hijacker
	ClientConnection net.Conn
	ClientReadWriter *bufio.ReadWriter
}

func (slaveResponseWriter *SlaveResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return slaveResponseWriter.ClientConnection, slaveResponseWriter.ClientReadWriter, nil
}

func CreateSlaveResponseWriter(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) *SlaveResponseWriter {
	slaveResponseWriter := new(SlaveResponseWriter)
	slaveResponseWriter.Body = new(bytes.Buffer)
	slaveResponseWriter.Code = 200
	slaveResponseWriter.ClientConnection = clientConnection
	slaveResponseWriter.ClientReadWriter = bufio.NewReadWriter(
		clientConnectionReader,
		clientConnectionWriter)
	return slaveResponseWriter
}

func CreateSlaveProxySession(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	args ...interface{}) {
	request, parsingError := http.ReadRequest(clientConnectionReader)
	if parsingError == nil {
		request.RemoteAddr = clientConnection.RemoteAddr().String()
		responseWriter := CreateSlaveResponseWriter(clientConnection, clientConnectionReader, clientConnectionWriter)
		args[0].(*goproxy.ProxyHttpServer).ServeHTTP(responseWriter, request)
		if responseWriter.Result().ContentLength > 0 {
			_ = responseWriter.Result().Write(clientConnectionWriter)
			_ = clientConnectionWriter.Flush()
		}
	} else {
		log.Print(parsingError)
		_ = clientConnection.Close()
	}
}

func StartHTTPSlave(address *string, port *string, useTLS *bool, proxy *goproxy.ProxyHttpServer) {
	log.Print("Starting HTTP server as slave")
	if *useTLS {
		log.Fatal("HTTPS proxy not supported yet")
	} else {
		MasterSlave.GeneralSlave(address, port, CreateSlaveProxySession, proxy)
	}
}

func StartHTTPMaster(address *string, port *string, useTLS *bool, proxy *goproxy.ProxyHttpServer) {
	log.Print("Starting HTTP server in Bind Mode")
	proxy.Verbose = true
	if *useTLS {
		log.Fatal("HTTPS proxy not supported yet")
	} else {
		listener := Sockets.Bind(address, port)
		if listener != nil {
			log.Fatal(http.Serve(listener, proxy))
		}
		log.Fatal("Something goes wrong when bind at: " + *address + ":" + *port)
	}
}

func StartHTTP(address *string, port *string, username *[]byte, password *[]byte, slave *bool, useTLS *bool) {
	passwordHash := Hashing.GetPasswordHashPasswordByteArray(username, password)
	proxy := goproxy.NewProxyHttpServer()
	if passwordHash != nil {
		proxy.OnRequest().DoFunc(func(request *http.Request, proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			authentication := strings.Split(request.Header.Get("Proxy-Authorization"), " ")
			if len(authentication) == 2 {
				if authentication[0] == "Basic" {
					rawCredentials, decodingError := base64.StdEncoding.DecodeString(authentication[1])
					if decodingError == nil {
						credentials := bytes.Split(rawCredentials, []byte(":"))
						if len(credentials) == 2 {
							if bytes.Equal(credentials[0], *username) && bytes.Equal(Hashing.PasswordHashingSHA3(credentials[1]), passwordHash) {
								log.Print("Login succeeded from: ", request.RemoteAddr)
								return request, nil
							}
						}
					}
				}
			}
			log.Print("Login failed with invalid credentials from: ", request.RemoteAddr)
			return request, goproxy.NewResponse(request, goproxy.ContentTypeText, http.StatusProxyAuthRequired, "Don't waste your time!")
		})
	}

	switch *slave {
	case true:
		StartHTTPSlave(address, port, useTLS, proxy)
	case false:
		StartHTTPMaster(address, port, useTLS, proxy)
	}
}
