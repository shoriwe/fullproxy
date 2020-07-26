package Basic

import (
	"FullProxy/FullProxy/Sockets"
	"bufio"
	"log"
	"net"
	"time"
)


func Proxy(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetConnectionReader *bufio.Reader,
	targetConnectionWriter *bufio.Writer,
	sniff bool) {
	for {
		for state := 0; state < 2; state++ {
			switch state {
			case 1:
				_ = clientConnection.SetReadDeadline(time.Now().Add(100))
				numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(clientConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				if sniff && numberOfBytesReceived > 0 {
					log.Print("Sending from: ", clientConnection.RemoteAddr(), " to: ", targetConnection.RemoteAddr(), " Chunk: ", buffer[:numberOfBytesReceived])
				}
				_, ConnectionError = Sockets.Send(targetConnectionWriter, buffer[:numberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			case 0:
				_ = targetConnection.SetReadDeadline(time.Now().Add(100))
				numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(targetConnectionReader, 20480)
				if ConnectionError != nil {
					if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
						return
					}
				}
				if sniff && numberOfBytesReceived > 0{
					log.Print("Sending from: ", targetConnection.RemoteAddr(), " to: ", clientConnection.RemoteAddr(), " Chunk: ", buffer[:numberOfBytesReceived])
				}
				_, ConnectionError = Sockets.Send(clientConnectionWriter, buffer[:numberOfBytesReceived])
				if ConnectionError != nil {
					return
				}
			}
		}
	}
}
