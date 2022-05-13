package http_hosts

import (
	"context"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
	"net/http"
	"net/url"
)

type Hosts struct {
	Client        *http.Client
	ListenAddress *net.TCPAddr
}

func (h *Hosts) Handle(_ net.Conn) error {
	panic("This should not be used")
}

func (h *Hosts) SetAuthenticationMethod(_ servers.AuthenticationMethod) {
}

func (h *Hosts) SetDial(dialFunc servers.DialFunc) {
	h.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return dialFunc(network, address)
			},
		},
	}
}

func (h *Hosts) SetListen(_ servers.ListenFunc) {
}

func (h *Hosts) SetListenAddress(address net.Addr) {
	h.ListenAddress = address.(*net.TCPAddr)
}

func (h *Hosts) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	host := request.URL.Hostname()
	if request.URL.Port() != "" {
		host = fmt.Sprintf("%s:%s", host, request.URL.Port())
	}
	u, parseError := url.Parse(fmt.Sprintf("%s://%s", request.URL.Scheme, host))
	if parseError != nil {
		return
	}
	u.Path = request.RequestURI
	targetRequest, newRequestError := http.NewRequest(request.Method, u.String(), request.Body)
	if newRequestError != nil {
		return
	}
	for key, values := range request.Header {
		for _, value := range values {
			targetRequest.Header.Add(key, value)
		}
	}
	response, requestError := h.Client.Do(targetRequest)
	if requestError != nil {
		return
	}
	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}
	writer.WriteHeader(response.StatusCode)
	_, copyError := io.Copy(writer, response.Body)
	if copyError != nil {
		return
	}
}

func NewHosts() servers.HTTPHandler {
	return &Hosts{
		Client: &http.Client{},
	}
}
