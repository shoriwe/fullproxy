package BindServer

import (
	"log"
	"net"
)


type Function func(conn net.Conn, username string, passwordHash []byte)


func BindServer(address string, port string, username string, passwordHash []byte, protocolFunction Function){
	server, BindingError  := net.Listen("tcp", address + ":" + port)
	if BindingError != nil {
		log.Print("Something goes wrong: " + BindingError.Error())
		return
	}
	for {
		clientConnection, _ := server.Accept()
		log.Print("Received connection from: ", clientConnection.RemoteAddr().String())
		go protocolFunction(clientConnection, username, passwordHash)
	}
}
