package Interface

import (
	"FullProxy/Sockets"
	"bufio"
	"log"
	"net"
)


type Function func(conn net.Conn, username string, passwordHash []byte)


func Client(address string, port string, username string, passwordHash []byte, function Function){
	log.Printf("Trying to connecto to %s:%s", address, port)
	var clientConnection net.Conn
	var masterConnection = Sockets.Connect(address,  port)
	if masterConnection == nil{
		return
	}
	masterConnectionReader := bufio.NewReader(masterConnection)
	log.Printf("Successfully connected to %s:%s", address, port)
	for {
		NumberOfReceivedBytes, buffer, ConnectionError := Sockets.Receive(masterConnectionReader, 1024)
		if ConnectionError != nil{
			log.Printf("Error: %s\n", ConnectionError.Error())
			break
		}
		if NumberOfReceivedBytes != 1{
			continue
		}
		if buffer[0] == Shutdown {
			log.Print("Received shutdown signal from the interface... Shuting down now the proxing service")
			break
		}
		clientConnection = Sockets.Connect(address, port)
		if clientConnection == nil{
			continue
		}
		go function(clientConnection, username, passwordHash)

	}
}