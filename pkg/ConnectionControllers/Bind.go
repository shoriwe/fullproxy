package ConnectionControllers

import (
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type Bind struct {
	Server        net.Listener
	ProxyProtocol ProxyProtocol
}

func (bind *Bind) Serve() error {
	for {
		clientConnection, connectionError := bind.Server.Accept()
		if connectionError != nil {
			return connectionError
		}
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go bind.ProxyProtocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
}
