package Sockets

import (
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"net"
)


func Send(connectionWriter ConnectionStructures.SocketWriter, data *[]byte) (int, error){
	var NumberOfBytesSent int
	var ConnectionError error
	NumberOfBytesSent, ConnectionError = connectionWriter.Write(*data)
	_ = connectionWriter.Flush()
	return NumberOfBytesSent, ConnectionError
}


func Receive(connectionReader ConnectionStructures.SocketReader, bufferSize int) (int, []byte, error){
	var receivedBytesError error
	buffer := make([]byte, bufferSize)
	NumberOfReceivedBytes, receivedBytesError := connectionReader.Read(buffer)
	if receivedBytesError != nil{
		return 0, nil, receivedBytesError
	}
	return NumberOfReceivedBytes, buffer, receivedBytesError
}


func Connect(address *string, port *string) net.Conn{
	var connection net.Conn
	var connectionError error
	connection, connectionError = net.Dial("tcp", *address + ":" + *port)
	if connectionError != nil{
		return nil
	}
	return connection
}
