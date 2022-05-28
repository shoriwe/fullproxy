package reverse

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	Target struct {
		RequestHeader  http.Header
		ResponseHeader http.Header
		URI            string
		CurrentHost    int
		Hosts          []*Host
	}
	HTTP struct {
		Targets                          map[string]*Target
		Dial                             servers.DialFunc
		WebSocketDialer                  *websocket.Dialer
		IncomingSniffer, OutgoingSniffer io.Writer
	}
)

func (H *HTTP) SetSniffers(incoming, outgoing io.Writer) {
	H.IncomingSniffer = incoming
	H.OutgoingSniffer = outgoing
}

func (target *Target) nextHost() *Host {
	if target.CurrentHost >= len(target.Hosts) {
		target.CurrentHost = 0
	}
	index := target.CurrentHost
	target.CurrentHost++
	return target.Hosts[index]
}

func (H *HTTP) SetListenAddress(_ net.Addr) {
}

func (H *HTTP) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (H *HTTP) SetListen(_ servers.ListenFunc) {
}

func (H *HTTP) Handle(_ net.Conn) error {
	panic("this should not be used")
}

func (H *HTTP) SetDial(dialFunc servers.DialFunc) {
	H.Dial = dialFunc
	H.WebSocketDialer = &websocket.Dialer{
		NetDial: dialFunc,
	}
}

func (H *HTTP) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target, found := H.Targets[request.Host]
	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if strings.Index(request.RequestURI, target.URI) != 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	host := target.nextHost()
	// Prepare new request
	u := &url.URL{
		Scheme: host.Scheme,
		Host:   host.Address,
		Path:   path.Join(host.URI, strings.Replace(request.RequestURI, target.URI, "/", 1)),
	}
	// Set request url
	newRequest, newRequestError := http.NewRequest(request.Method, u.String(), request.Body)
	if newRequestError != nil {
		// TODO: Do something with the error
		return
	}
	newRequest.Header = request.Header.Clone()
	// Inject Headers in request
	for key, values := range target.RequestHeader {
		newRequest.Header[key] = values
	}
	newRequest.Host = newRequest.Header.Get("Host")
	// Check websocket
	if websocket.IsWebSocketUpgrade(newRequest) {
		if H.OutgoingSniffer != nil {
			sniffError := newRequest.Write(H.OutgoingSniffer)
			if sniffError != nil {
				// TODO: Do something with the error
				return
			}
			_, _ = fmt.Fprintf(H.OutgoingSniffer, common.SniffSeparator)
		}
		switch host.Scheme {
		case "http", "ws":
			u.Scheme = "ws"
		case "https", "wss":
			u.Scheme = "wss"
		}
		newRequest.Header.Del("Upgrade")
		newRequest.Header.Del("Connection")
		newRequest.Header.Del("Sec-Websocket-Key")
		newRequest.Header.Del("Sec-Websocket-Version")
		newRequest.Header.Del("Sec-Websocket-Extensions")
		targetConnection, response, dialError := H.WebSocketDialer.Dial(
			u.String(),
			newRequest.Header,
		)
		if H.IncomingSniffer != nil {
			sniffError := response.Write(H.IncomingSniffer)
			if sniffError != nil {
				// TODO: Do something with the error
				return
			}
			_, _ = fmt.Fprintf(H.IncomingSniffer, common.SniffSeparator)
		}
		if dialError != nil {
			// TODO: Do something with the error
			return
		}
		// Upgrade connection
		clientResponseHeader := response.Header.Clone()
		for key, values := range target.ResponseHeader {
			clientResponseHeader[key] = values
		}
		upgrader := &websocket.Upgrader{
			ReadBufferSize:  host.WebsocketReadBufferSize,
			WriteBufferSize: host.WebsocketWriteBufferSize,
		}
		clientConnection, upgradeError := upgrader.Upgrade(
			writer,
			request,
			clientResponseHeader,
		)
		if upgradeError != nil {
			// TODO: Do something with the error
			return
		}

		forwardError := common.ForwardWebsocketsTraffic(clientConnection, targetConnection, H.IncomingSniffer, H.OutgoingSniffer)
		if forwardError != nil {
			// TODO: Do something with the error
		}
		return
	}
	// Prepare client
	serverConnection, connectionError := H.Dial(host.Network, host.Address)
	if connectionError != nil {
		// TODO: Do something with the error
		return
	}
	defer serverConnection.Close()
	if host.TLSConfig != nil {
		serverConnection = tls.Client(serverConnection, host.TLSConfig)
	}
	server := &common.Sniffer{
		WriteSniffer: H.IncomingSniffer,
		ReadSniffer:  H.OutgoingSniffer,
		Connection:   serverConnection,
	}
	// Send request to server
	sendRequestError := newRequest.Write(server)
	if sendRequestError != nil {
		// TODO: Do something with the error
		return
	}
	// Receive server response
	serverResponse, readResponseError := http.ReadResponse(bufio.NewReader(server), newRequest)
	if readResponseError != nil {
		// TODO: Do something with the error
		return
	}
	defer serverResponse.Body.Close()
	// Inject response headers
	for key, values := range target.ResponseHeader {
		writer.Header()[key] = values
	}
	writer.WriteHeader(serverResponse.StatusCode)
	_, copyError := io.Copy(writer, serverResponse.Body)
	if copyError != nil {
		// TODO: Do something with the error
		return
	}
}

func NewHTTP(targets map[string]*Target) servers.HTTPHandler {
	return &HTTP{
		Targets: targets,
	}
}
