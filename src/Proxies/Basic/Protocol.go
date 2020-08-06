package Basic

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"log"
	"net"
	"time"
)


func HandleReadWrite(
	sourceConnection net.Conn, destinationAddress string,
	sourceReader ConnectionStructures.SocketReader, destinationWriter ConnectionStructures.SocketWriter,
	connectionAlive *bool){
	for tries := 0; tries < 5; tries++{
		if !(*connectionAlive){
			break
		}
		_ = sourceConnection.SetReadDeadline(time.Now().Add(20 * time.Second))
		numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(sourceReader, 20480)
		if ConnectionError != nil {
			if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
				break
			}
		} else {
			tries = 0
		}
		if numberOfBytesReceived > 0 {
			realChunk := buffer[:numberOfBytesReceived]
			_, ConnectionError = Sockets.Send(destinationWriter, &realChunk)
			if ConnectionError != nil {
				break
			}
			log.Print("Sending: ", realChunk, " From: ", sourceConnection.RemoteAddr(), " To: ", destinationAddress)
			realChunk = nil
		}
		buffer = nil
	}
	_ = sourceConnection.Close()
	*connectionAlive = false
}


func Proxy(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader ConnectionStructures.SocketReader,
	clientConnectionWriter ConnectionStructures.SocketWriter,
	targetConnectionReader ConnectionStructures.SocketReader,
	targetConnectionWriter ConnectionStructures.SocketWriter) {
	connectionAlive := true
	go HandleReadWrite(clientConnection, targetConnection.RemoteAddr().String(), clientConnectionReader, targetConnectionWriter, &connectionAlive)
	go HandleReadWrite(targetConnection, clientConnection.RemoteAddr().String(), targetConnectionReader, clientConnectionWriter, &connectionAlive)

}
