package SetupControllers

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
)

func Bind(host *string, port *string, proxy ConnectionControllers.ProxyProtocol) {
	server, bindError := Sockets.Bind(host, port)
	if bindError != nil {
		log.Fatal("Could not bind to wanted address")
	}
	log.Print("Bind successfully")
	controller := new(ConnectionControllers.Bind)
	controller.Server = server
	controller.ProxyProtocol = proxy
	log.Fatal(controller.Serve())
}
