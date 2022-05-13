package listeners

import (
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"log"
	"net"
	"net/http"
)

type listenerInboundFilter struct {
	Listener Listener
}

func (l *listenerInboundFilter) Close() error {
	return l.Listener.Close()
}

func (l *listenerInboundFilter) Addr() net.Addr {
	return l.Listener.Addr()
}

func (l *listenerInboundFilter) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	filterError := l.Listener.Filter().Inbound(conn.RemoteAddr().String())
	if filterError != nil {
		_ = conn.Close()
		return nil, filterError
	}
	return conn, nil
}

func ServeHTTPHandler(listener Listener, handler servers.HTTPHandler, logFunc LogFunc) error {
	handler.SetDial(func(network, address string) (net.Conn, error) {
		filterError := listener.Filter().Outbound(address)
		if filterError != nil {
			return nil, filterError
		}
		return listener.Dial(network, address)
	})
	handler.SetListen(func(network, address string) (net.Listener, error) {
		filterError := listener.Filter().Listen(address)
		if filterError != nil {
			return nil, filterError
		}
		l, listenError := listener.Listen(network, address)
		if listenError != nil {
			return nil, listenError
		}
		return &TCPListener{l.(*net.TCPListener), listener.Filter()}, nil
	})
	handler.SetListenAddress(listener.Addr())
	initError := listener.Init()
	if initError != nil {
		return initError
	}
	if logFunc == nil {
		logFunc = func(args ...interface{}) {
			log.Println(args)
		}
	}
	server := http.Server{
		Addr:    listener.Addr().String(),
		Handler: handler,
	}
	return server.Serve(&listenerInboundFilter{Listener: listener})
}
