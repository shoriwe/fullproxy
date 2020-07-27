package Interface

import (
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"log"
	"net"
)


type Function func(conn net.Conn, username *[]byte, passwordHash *[]byte)


func Client(address string, port string, username *[]byte, passwordHash *[]byte, function Function){
	log.Printf("Trying to connecto to %s:%s", address, port)
	masterConnection := Sockets.Connect(address,  port)
	if masterConnection != nil{
		masterConnectionReader := bufio.NewReader(masterConnection)
		log.Printf("Successfully connected to %s:%s", address, port)
		for {
			NumberOfReceivedBytes, buffer, ConnectionError := Sockets.Receive(masterConnectionReader, 1024)
			if ConnectionError == nil{
				if NumberOfReceivedBytes == 1{
					switch buffer[0] {
					case Shutdown:
						log.Print("Received shutdown signal... shutting down")
						break
					case NewConnection:
						clientConnection := Sockets.Connect(address, port)
						if clientConnection != nil{
							go function(clientConnection, username, passwordHash)
						}
					}
				} else {
					continue
				}
			} else {
				log.Printf("Error: %s\n", ConnectionError.Error())
				break
			}
		}
	}
}