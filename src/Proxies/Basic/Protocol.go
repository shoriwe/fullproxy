package Basic

import (
	"bufio"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
	"time"
)


func HandleReadWrite(
	sourceConnection net.Conn, sourceReader *bufio.Reader, destinationWriter *bufio.Writer, connectionAlive *bool){
	for {
		if !(*connectionAlive){
			break
		}
		_ = sourceConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
		numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(sourceReader, 20480)
		if ConnectionError != nil {
			if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
				break
			}
		}
		if numberOfBytesReceived > 0 {
			_, ConnectionError = Sockets.Send(destinationWriter, buffer[:numberOfBytesReceived])
			if ConnectionError != nil {
				break
			}}
		buffer = nil
	}
	_ = sourceConnection.Close()
	*connectionAlive = false
}


func Proxy(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetConnectionReader *bufio.Reader,
	targetConnectionWriter *bufio.Writer) {
	connectionAlive := true

	go HandleReadWrite(clientConnection, clientConnectionReader, targetConnectionWriter, &connectionAlive)
	go HandleReadWrite(targetConnection, targetConnectionReader, clientConnectionWriter, &connectionAlive)
}
