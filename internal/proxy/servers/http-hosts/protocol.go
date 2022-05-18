package http_hosts

import (
	"context"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"io"
	"net"
	"net/http"
)

type Hosts struct {
	Client                           *http.Client
	ListenAddress                    *net.TCPAddr
	IncomingSniffer, OutgoingSniffer io.Writer
}

func (h *Hosts) SetSniffers(incoming, outgoing io.Writer) {
	h.IncomingSniffer = incoming
	h.OutgoingSniffer = outgoing
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
	_ = common.SniffRequest(h.IncomingSniffer, request)
	targetRequest, newRequestError := http.NewRequest(request.Method, request.URL.String(), request.Body)
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
	_ = common.SniffResponse(h.OutgoingSniffer, response)
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
