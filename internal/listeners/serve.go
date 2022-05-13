package listeners

import (
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"log"
	"net"
)

func Serve(listener Listener, protocol servers.Protocol, logFunc LogFunc) error {
	if listener.Filter() == nil {
		listener.SetFilters(&NoFilter{})
	}
	protocol.SetDial(func(network, address string) (net.Conn, error) {
		filterError := listener.Filter().Outbound(address)
		if filterError != nil {
			return nil, filterError
		}
		return listener.Dial(network, address)
	})
	protocol.SetListen(func(network, address string) (net.Listener, error) {
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
	protocol.SetListenAddress(listener.Addr())
	initError := listener.Init()
	if initError != nil {
		return initError
	}
	if logFunc == nil {
		logFunc = func(args ...interface{}) {
			log.Println(args)
		}
	}
	for {
		clientConnection, connectionError := listener.Accept()
		if connectionError != nil {
			return connectionError
		}
		go func() {
			logFunc(fmt.Sprintf("Connection from: %s", clientConnection.RemoteAddr().String()))
			defer clientConnection.Close()
			filterError := listener.Filter().Inbound(clientConnection.RemoteAddr().String())
			if filterError != nil {
				logFunc(filterError)
				return
			}
			handleError := protocol.Handle(clientConnection)
			if handleError != nil {
				logFunc(handleError)
				return
			}
		}()
	}
}
