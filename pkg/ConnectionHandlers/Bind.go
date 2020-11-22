package ConnectionHandlers

import (
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"net"
)

func Bind(address *string, port *string, protocol ProxyProtocol) {
	server, BindingError := net.Listen("tcp", *address+":"+*port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	log.Printf("Listening on %s:%s", *address, *port)
	for {
		clientConnection, _ := server.Accept()
		log.Print("Received connection from: ", clientConnection.RemoteAddr().String())
		clientConnectionReader, clientConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(clientConnection)
		go protocol.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
	}
}