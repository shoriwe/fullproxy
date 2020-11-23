package ConnectionControllers

import (
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type Bind struct {
	Server        net.Listener
	ProxyProtocol ProxyProtocol
	LoggingMethod LoggingMethod
}

func (bind *Bind) SetLoggingMethod(loggingMethod LoggingMethod) error {
	bind.LoggingMethod = loggingMethod
	return nil
}

func (bind *Bind) Serve() error {
	for {
		clientConnection, connectionError := bind.Server.Accept()
		LogData(bind.LoggingMethod, "Client connection received from: ", clientConnection.RemoteAddr().String())
		if connectionError != nil {
			LogData(bind.LoggingMethod, connectionError)
			return connectionError
		}
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go bind.ProxyProtocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
}
