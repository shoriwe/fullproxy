package MasterSlave

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
)


type Function func(conn net.Conn, connReader ConnectionStructures.SocketReader, connWriter ConnectionStructures.SocketWriter, args...interface{})


func Slave(address *string, port *string, function Function, args...interface{}){
	log.Printf("Trying to connecto to %s:%s", *address, *port)
	masterConnection := Sockets.Connect(address,  port)
	if masterConnection != nil{
		masterConnectionReader, _ := ConnectionStructures.CreateReaderWriter(masterConnection)
		log.Printf("Successfully connected to %s:%s", *address, *port)
		for {
			NumberOfReceivedBytes, buffer, ConnectionError := Sockets.Receive(masterConnectionReader, 1024)
			if ConnectionError == nil{
				if NumberOfReceivedBytes == 1{
					switch buffer[0] {
					case Shutdown[0]:
						log.Print("Received shutdown signal... shutting down")
						break
					case NewConnection[0]:
						clientConnection := Sockets.Connect(address, port)
						if clientConnection != nil{
							clientConnectionReader, clientConnectionWriter := ConnectionStructures.CreateTunnelReaderWriter(clientConnection)
							if clientConnectionReader != nil && clientConnectionWriter != nil {
								go function(clientConnection, clientConnectionReader, clientConnectionWriter, args...)
							} else {
								_ = clientConnection.Close()
							}
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