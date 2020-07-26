package Sockets

import (
	"bufio"
	"net"
)


func Send(connectionWriter *bufio.Writer, data []byte) (int, error){
	NumberOfBytesSent, ConnectionError := connectionWriter.Write(data)
	_ = connectionWriter.Flush()
	return NumberOfBytesSent, ConnectionError
}


func Receive(connectionReader *bufio.Reader, bufferSize int) (int, []byte, error){
	buffer := make([]byte, bufferSize)
	NumberOfReceivedBytes, receivedBytesError := connectionReader.Read(buffer)
	if receivedBytesError != nil{
		return 0, nil, receivedBytesError
	}
	return NumberOfReceivedBytes, buffer, receivedBytesError
}


func Connect(ip string, port string) net.Conn{
	var connection net.Conn
	var connectionError error
	connection, connectionError = net.Dial("tcp", ip + ":" + port)
	if connectionError != nil{
		return nil
	}
	return connection
}
