package ConnectionStructures

import (
	"bufio"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"log"
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
	var finalKey []byte
	myKeyPart := GenerateKey(8)
	myRSAKey, keyGenerationError := rsa.GenerateKey(rand.Reader, 512)
	if keyGenerationError == nil {
		publicKeyBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&myRSAKey.PublicKey)})
		publicKeyLength := make([]byte, 4)
		binary.LittleEndian.PutUint32(publicKeyLength, uint32(len(publicKeyBytes)))
		_, connectionError := sourceWriter.Write(append(publicKeyLength, publicKeyBytes...))
		_ = sourceWriter.Flush()
		if connectionError == nil {
			lengthOfOtherPublicKey := make([]byte, 4)
			numberOfBytesReceived, connectionError := sourceReader.Read(lengthOfOtherPublicKey)
			if connectionError == nil && numberOfBytesReceived == 4{
				otherPublicKeyBytes := make([]byte, binary.LittleEndian.Uint32(lengthOfOtherPublicKey))
				numberOfBytesReceived, connectionError := sourceReader.Read(otherPublicKeyBytes)
				if connectionError == nil && numberOfBytesReceived > 0 {
					block, _ := pem.Decode(otherPublicKeyBytes[:numberOfBytesReceived])
					otherPublicKey, parsingError := x509.ParsePKCS1PublicKey(block.Bytes)
					if parsingError == nil {
						myEncryptedKeyPart, encryptionError := rsa.EncryptPKCS1v15(rand.Reader, otherPublicKey, myKeyPart)
						if encryptionError == nil {
							_, connectionError := sourceWriter.Write(myEncryptedKeyPart)
							_ = sourceWriter.Flush()
							if connectionError == nil {
								otherEncryptedKeyPart := make([]byte, 1024)
								numberOfBytesReceived, connectionError := sourceReader.Read(otherEncryptedKeyPart)
								if connectionError == nil && numberOfBytesReceived > 0 {
									otherKeyPart, decryptionError := rsa.DecryptPKCS1v15(rand.Reader, myRSAKey, otherEncryptedKeyPart[:numberOfBytesReceived])
									if decryptionError == nil {
										if len(otherKeyPart) > 0 {
											if myKeyPart[0] > otherKeyPart[0] {
												finalKey = append(myKeyPart, otherKeyPart...)
											} else {
												finalKey = append(otherKeyPart, myKeyPart...)
											}
										} else {
											log.Print("No key received")
										}
									} else {
										log.Print(decryptionError)
									}
								} else {
									log.Print(connectionError)
									log.Print("Received just: ", numberOfBytesReceived)
								}
							} else {
								log.Print(connectionError)
							}
						} else {
							log.Print(encryptionError)
						}
					} else {
						log.Print(parsingError)
					}
				} else {
					log.Print(connectionError)
					log.Print("Received just: ", numberOfBytesReceived)
				}
			}else {
				log.Println(connectionError)
				log.Println(numberOfBytesReceived)
			}
		} else {
			log.Print(connectionError)
		}
	} else {
		log.Print(keyGenerationError)
	}
	return finalKey
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
		if len(key) == 16 {
			cipherBlockReader, _ := aes.NewCipher(key)
			cipherBlockWriter, _ := aes.NewCipher(key)
			tunnelReader.CipherBlock = cipherBlockReader
			tunnelWriter.CipherBlock = cipherBlockWriter
			return tunnelReader, tunnelWriter
		}
	}
	return nil, nil
}
