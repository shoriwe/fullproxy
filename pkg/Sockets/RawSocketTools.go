package Sockets

import (
	"bufio"
	"net"
)

func Send(connectionWriter *bufio.Writer, data *[]byte) (int, error) {
	var NumberOfBytesSent int
	var ConnectionError error
	NumberOfBytesSent, ConnectionError = connectionWriter.Write(*data)
	_ = connectionWriter.Flush()
	return NumberOfBytesSent, ConnectionError
}

func Receive(connectionReader *bufio.Reader, bufferSize int) (int, []byte, error) {
	var receivedBytesError error
	buffer := make([]byte, bufferSize)
	NumberOfReceivedBytes, receivedBytesError := connectionReader.Read(buffer)
	if receivedBytesError != nil {
		return 0, nil, receivedBytesError
	}
	return NumberOfReceivedBytes, buffer, receivedBytesError
}

func Connect(host *string, port *string) (net.Conn, error) {
	return net.Dial("tcp", *host+":"+*port)
}

func Bind(host *string, port *string) (net.Listener, error) {
	return net.Listen("tcp", *host+":"+*port)
}
