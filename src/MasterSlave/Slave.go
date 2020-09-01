package MasterSlave

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
	"time"
)


type Function func(conn net.Conn, connReader ConnectionStructures.SocketReader, connWriter ConnectionStructures.SocketWriter, args...interface{})


func GeneralSlave(masterAddress *string, masterPort *string, function Function, args...interface{}){
	log.Printf("Trying to connecto to %s:%s", *masterAddress, *masterPort)
	masterConnection := Sockets.Connect(masterAddress, masterPort)
	if masterConnection != nil{
		masterConnectionReader, _ := ConnectionStructures.CreateReaderWriter(masterConnection)
		log.Printf("Successfully connected to %s:%s", *masterAddress, *masterPort)
		for {
			_ = masterConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
			NumberOfReceivedBytes, buffer, connectionError := Sockets.Receive(masterConnectionReader, 1024)
			if connectionError == nil{
				if NumberOfReceivedBytes == 1{
					switch buffer[0] {
					case Shutdown[0]:
						log.Print("Received shutdown signal... shutting down")
						break
					case NewConnection[0]:
						clientConnection := Sockets.Connect(masterAddress, masterPort)
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
			} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				break
			}
		}
	}
}


func RemotePortForwardSlave(masterAddress *string, masterPort *string, localAddress *string, localPort *string){
	localServer := Sockets.Bind(localAddress, localPort)
	masterConnection := Sockets.Connect(masterAddress, masterPort)
	masterConnectionReader, masterConnectionWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(masterConnection)
	if masterConnection != nil {
		for {
			clientConnection, connectionError := localServer.Accept()
			if connectionError == nil {
				_, connectionError := Sockets.Send(masterConnectionWriter, &NewConnection)
				if connectionError == nil {
					_ = masterConnection.SetReadDeadline(time.Now().Add(3 * time.Second))
					numberOfBytesReceived, response, connectionError := Sockets.Receive(masterConnectionReader, 1)
					if connectionError == nil {
						if numberOfBytesReceived == 1 {
							switch response[0] {
							case NewConnection[0]:
								targetConnection := Sockets.Connect(masterAddress, masterPort)
								if targetConnection != nil {
									go startGeneralProxying(clientConnection, targetConnection)
								} else {
									_ = clientConnection.Close()
									log.Fatal("Connectivity issues with master server")
								}
							case FailToConnectToTarget[0]:
								break
							case UnknownOperation[0]:
								break
							}
						}
					} else if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
						break
					}
				}
			}
		}
	}
}