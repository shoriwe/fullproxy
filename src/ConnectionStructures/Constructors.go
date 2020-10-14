package ConnectionStructures

import (
	"bufio"
	"net"
)

func createReaderWriter(connection net.Conn) (*bufio.Reader, *bufio.Writer) {
	return bufio.NewReader(connection), bufio.NewWriter(connection)
}

func CreateSocketConnectionReaderWriter(connection net.Conn) (SocketReader, SocketWriter) {
	var socketConnectionReader BasicConnectionReader
	var socketConnectionWriter BasicConnectionWriter
	socketConnectionReader.Reader, socketConnectionWriter.Writer = createReaderWriter(connection)
	return &socketConnectionReader, &socketConnectionWriter
}
