package BindServer

import (
	"log"
	"net"
)


type Function func(conn net.Conn, args...interface{})


func Bind(address *string, port *string, protocolFunction Function, args...interface{}){
	server, BindingError  := net.Listen("tcp", *address + ":" + *port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	for {
		clientConnection, _ := server.Accept()
		log.Print("Received connection from: ", clientConnection.RemoteAddr().String())
		go protocolFunction(clientConnection, args[:]...)
	}
}
