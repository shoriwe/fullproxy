package proxies

import (
	"net"
	"net/http"

	"github.com/elazarl/goproxy"
)

type HTTP struct {
	Listener net.Listener
	// TODO: Authentication configuration
	Dial func(network, addr string) (net.Conn, error)
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
