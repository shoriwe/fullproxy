package Basic

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
	"time"
)


func HandleReadWrite(
	sourceConnection net.Conn, destinationAddress string,
	sourceReader ConnectionStructures.SocketReader, destinationWriter ConnectionStructures.SocketWriter,
	connectionAlive *bool){
	for tries := 0; tries < 5; tries++{
		_ = sourceConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(sourceReader, 20480)
		if connectionError != nil {
			// If the error is not "Timeout"
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				break
			} else {
				if !(*connectionAlive){
					break
				}
			}
		} else {
			tries = 0
		}
		if numberOfBytesReceived > 0 {
			realChunk := buffer[:numberOfBytesReceived]
			_, connectionError = Sockets.Send(destinationWriter, &realChunk)
			if connectionError != nil {
				break
			}
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
