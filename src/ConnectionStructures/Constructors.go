package ConnectionStructures

import (
	"bufio"
	"crypto/aes"
	"net"
)


func NegotiateKey(
	sourceReader SocketReader,
	sourceWriter SocketWriter) []byte {
	key := []byte("1234567890123456")

	return key
}


func CreateReaderWriter(connection net.Conn) (*bufio.Reader, *bufio.Writer){
	return bufio.NewReader(connection), bufio.NewWriter(connection)
}


func CreateSocketConnectionReaderWriter(connection net.Conn) (SocketReader, SocketWriter){
	var socketConnectionReader BasicConnectionReader
	var socketConnectionWriter BasicConnectionWriter
	socketConnectionReader.Reader, socketConnectionWriter.Writer = CreateReaderWriter(connection)
	return &socketConnectionReader, &socketConnectionWriter
}


func CreateTunnelReaderWriter(connection net.Conn) (SocketReader, SocketWriter){
	var tunnelReader TunnelReader
	var tunnelWriter TunnelWriter
	tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter = CreateSocketConnectionReaderWriter(connection)
	cipherBlockReader, _ := aes.NewCipher(NegotiateKey(tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter))
	cipherBlockWriter, _ := aes.NewCipher(NegotiateKey(tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter))
	tunnelReader.CipherBlock = cipherBlockReader
	tunnelWriter.CipherBlock = cipherBlockWriter
	return &tunnelReader, &tunnelWriter
}
