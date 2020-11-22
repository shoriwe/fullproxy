package BindServer

import (
	"bufio"
	"github.com/shoriwe/FullProxy/pkg/ConnectionStructures"
	"log"
	"net"
)

type Function func(
	conn net.Conn,
	connReader *bufio.Reader,
	connWriter *bufio.Writer,
	args ...interface{})

func Bind(address *string, port *string, protocolFunction Function, args ...interface{}) {
	server, BindingError := net.Listen("tcp", *address+":"+*port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	log.Printf("Listening on %s:%s", *address, *port)
	for {
		clientConnection, _ := server.Accept()
		log.Print("Received connection from: ", clientConnection.RemoteAddr().String())
		clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(clientConnection)
		go protocolFunction(clientConnection, clientConnectionReader, clientConnectionWriter, args[:]...)
	}
}
