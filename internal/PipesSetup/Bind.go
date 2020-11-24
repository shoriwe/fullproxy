package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"log"
)

func Bind(host *string, port *string, proxy Types.ProxyProtocol) {
	server, bindError := Sockets.Bind(host, port)
	if bindError != nil {
		log.Fatal("Could not bind to wanted address")
	}
	log.Print("Bind successfully in: ", *host, ":", *port)
	pipe := new(Pipes.Bind)
	_ = pipe.SetLoggingMethod(log.Print)
	pipe.Server = server
	pipe.ProxyProtocol = proxy
	Pipes.Serve(pipe)
}
