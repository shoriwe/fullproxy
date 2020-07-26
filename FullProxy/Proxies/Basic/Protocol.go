package Basic

import (
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"log"
	"net"
)


func HandleReadWrite(
	sourceConnection net.Conn, destinationConnection net.Conn,
	sourceReader *bufio.Reader, destinationWriter *bufio.Writer, sniff *bool){
	for {
		// _ = clientConnection.SetReadDeadline(time.Now().Add(100))
		numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(sourceReader, 20480)
		if ConnectionError != nil {
			if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
				break
			}
		}
		if *sniff && numberOfBytesReceived > 0 {
			log.Print("Sending from: ", sourceConnection.RemoteAddr(), " to: ", destinationConnection.RemoteAddr(), " Chunk: ", buffer[:numberOfBytesReceived])
		}
		_, ConnectionError = Sockets.Send(destinationWriter, buffer[:numberOfBytesReceived])
		if ConnectionError != nil {
			break
		}
		buffer = nil
	}
	_ = sourceConnection.Close()
}

func Proxy(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetConnectionReader *bufio.Reader,
	targetConnectionWriter *bufio.Writer,
	sniff bool) {
	go HandleReadWrite(clientConnection, targetConnection, clientConnectionReader, targetConnectionWriter, &sniff)
	go HandleReadWrite(targetConnection, clientConnection, targetConnectionReader, clientConnectionWriter, &sniff)
}
