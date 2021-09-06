package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"log"
	"net"
)

func Bind(networkType string, address string, proxy Types.ProxyProtocol, inboundFilter Types.IOFilter) {
	server, bindError := net.Listen(networkType, address)
	if bindError != nil {
		log.Fatal("Could not bind to wanted address")
	}
	log.Print("Bind successfully in: ", address)
	pipe := &Pipes.Bind{
		Server:        server,
		ProxyProtocol: proxy,
		LoggingMethod: log.Print,
		InboundFilter: inboundFilter,
	}
	Pipes.Serve(pipe)
}
