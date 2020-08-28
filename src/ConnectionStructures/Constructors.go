package ConnectionStructures

import (
	"bufio"
	"crypto/aes"
	"crypto/rand"
	"net"
)


func GenerateKey(keyLength int) []byte{
	key := make([]byte, keyLength)
	_, _ = rand.Reader.Read(key)
	return key
}


func NegotiateKey(
	sourceReader SocketReader,
	sourceWriter SocketWriter) []byte {
	return []byte("passwordpassword")
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
	tunnelReader := new(TunnelReader)
	tunnelWriter := new(TunnelWriter)
	tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter = CreateSocketConnectionReaderWriter(connection)
	key := NegotiateKey(tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter)
	if key != nil {
		cipherBlockReader, _ := aes.NewCipher(key)
		cipherBlockWriter, _ := aes.NewCipher(key)
		tunnelReader.CipherBlock = cipherBlockReader
		tunnelWriter.CipherBlock = cipherBlockWriter
		return tunnelReader, tunnelWriter
	}
	return nil, nil
}
