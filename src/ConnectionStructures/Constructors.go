package ConnectionStructures

import (
	"bufio"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net"
)



// Function Extracted from https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		log.Print("is encrypted pem block")
		b, _ = x509.DecryptPEMBlock(block, nil)
	}
	ifc, _ := x509.ParsePKIXPublicKey(b)
	key, _ := ifc.(*rsa.PublicKey)
	return key
}


// Function Extracted from https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, _ := x509.MarshalPKIXPublicKey(pub)

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}


func NegotiateKey(
	sourceReader SocketReader,
	sourceWriter SocketWriter) []byte {

	myKey := make([]byte, 16)
	_, _ = rand.Read(myKey)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 512)
	publicKeyBytes := PublicKeyToBytes(&privateKey.PublicKey)

	_, _ = sourceWriter.Write(publicKeyBytes)
	_ = sourceWriter.Flush()

	otherPublicKeyBytes := make([]byte, 2048)
	numberOfReceivedBytes, _ := sourceReader.Read(otherPublicKeyBytes)
	otherPublicKey := BytesToPublicKey(otherPublicKeyBytes[:numberOfReceivedBytes])

	myPasswordChunk, _ := rsa.EncryptPKCS1v15(rand.Reader, otherPublicKey, myKey)

	_, _ = sourceWriter.Write(myPasswordChunk)
	_ = sourceWriter.Flush()

	otherRawPasswordChunkEncrypted := make([]byte, 2048)
	numberOfReceivedBytes, _ = sourceReader.Read(otherRawPasswordChunkEncrypted)
	otherPasswordChunkEncrypted := otherRawPasswordChunkEncrypted[:numberOfReceivedBytes]

	otherPasswordChunk, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, otherPasswordChunkEncrypted)

	var key []byte
	if otherPasswordChunk[0] > myKey[0]{
		key = append(myKey, otherPasswordChunk...)
	} else {
		key = append(otherPasswordChunk, myKey...)
	}
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
	key := NegotiateKey(tunnelReader.ActiveSocketReader, tunnelWriter.ActiveSocketWriter)
	cipherBlockReader, _ := aes.NewCipher(key)
	cipherBlockWriter, _ := aes.NewCipher(key)
	tunnelReader.CipherBlock = cipherBlockReader
	tunnelWriter.CipherBlock = cipherBlockWriter
	return &tunnelReader, &tunnelWriter
}
