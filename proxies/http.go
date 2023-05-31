package proxies

import (
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

type HTTP struct {
	Listener net.Listener
	// TODO: Authentication configuration
	Dial network.DialFunc
}

func (h *HTTP) Close() {
	h.Listener.Close()
}

func (h *HTTP) Addr() net.Addr {
	return h.Listener.Addr()
}

func (h *HTTP) Serve() error {
	server := goproxy.NewProxyHttpServer()
	server.ConnectDial = h.Dial
	return http.Serve(h.Listener, server)
}
