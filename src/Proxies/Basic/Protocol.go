package Basic

import (
	"bufio"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
	"time"
)

func HandleReadWrite(
	sourceConnection net.Conn,
	sourceReader *bufio.Reader,
	destinationWriter *bufio.Writer,
	connectionAlive *bool) {
	for tries := 0; tries < 5; tries++ {
		_ = sourceConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
		numberOfBytesReceived, buffer, connectionError := Sockets.Receive(sourceReader, 1048576)
		if connectionError != nil {
			// If the error is not "Timeout"
			if parsedConnectionError, ok := connectionError.(net.Error); !(ok && parsedConnectionError.Timeout()) {
				break
			} else {
				if !(*connectionAlive) {
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
	if *connectionAlive {
		*connectionAlive = false
	}
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
