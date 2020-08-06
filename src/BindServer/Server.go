package BindServer

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"log"
	"net"
)


type Function func(conn net.Conn, connReader  ConnectionStructures.SocketReader, connWriter ConnectionStructures.SocketWriter, args...interface{})


func Bind(address *string, port *string, protocolFunction Function, args...interface{}){
	server, BindingError  := net.Listen("tcp", *address + ":" + *port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	for {
		clientConnection, _ := server.Accept()
		log.Print("Received connection from: ", clientConnection.RemoteAddr().String())
		clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateReaderWriter(clientConnection)
		go protocolFunction(clientConnection, clientConnectionReader, clientConnectionWriter, args[:]...)
	}
}
