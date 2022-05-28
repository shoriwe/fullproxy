package reverse

import (
	"context"
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
	defer request.Body.Close()
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
		dialer := &websocket.Dialer{
			NetDial: func(_, _ string) (net.Conn, error) {
				return H.Dial(host.Network, host.Address)
			},
			TLSClientConfig: host.TLSConfig,
		}
		targetConnection, response, dialError := dialer.Dial(
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
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return H.Dial(host.Network, host.Address)
			},
			TLSClientConfig: host.TLSConfig,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	newRequest.Body = &common.RequestSniffer{
		HeaderDone: false,
		Writer:     H.OutgoingSniffer,
		Request:    request,
	}
	response, requestError := client.Do(newRequest)
	if requestError != nil {
		// TODO: Do something with the error
		return
	}
	defer response.Body.Close()
	newResponse := &http.Response{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Proto:      response.Proto,
		ProtoMajor: response.ProtoMajor,
		ProtoMinor: response.ProtoMinor,
		Header:     response.Header.Clone(),
		Body: &common.ResponseSniffer{
			HeaderDone: false,
			Writer:     H.IncomingSniffer,
			Response:   response,
		},
		ContentLength:    response.ContentLength,
		TransferEncoding: response.TransferEncoding,
		Close:            response.Close,
		Uncompressed:     response.Uncompressed,
		Trailer:          response.Trailer.Clone(),
		Request:          response.Request,
		TLS:              response.TLS,
	}
	for key, values := range newResponse.Header {
		writer.Header()[key] = values
	}
	for key, values := range target.ResponseHeader {
		writer.Header()[key] = values
	}
	writer.WriteHeader(newResponse.StatusCode)
	_, copyError := io.Copy(writer, newResponse.Body)
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
