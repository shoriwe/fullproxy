package ConnectionStructures

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func GenerateIV() ([]byte, error) {
	iv := make([]byte, 16)
	_, randomReadError := rand.Read(iv)
	if randomReadError == nil {
		return iv, nil
	}
	return nil, randomReadError
}

type TunnelReader struct {
	CipherBlock        cipher.Block
	ActiveSocketReader SocketReader
}

func (tunnelReader *TunnelReader) Read(buffer []byte) (int, error) {
	tempBuffer := make([]byte, len(buffer)+16)
	numberOfReceivedBytes, connectionError := tunnelReader.ActiveSocketReader.Read(tempBuffer)
	if connectionError == nil {
		if numberOfReceivedBytes >= 16 {
			iv := tempBuffer[:16]
			decryptCBCBlock := cipher.NewCBCDecrypter(tunnelReader.CipherBlock, iv)
			decryptedBufferPadded := make([]byte, 0)
			lengthEncryptedMessage := numberOfReceivedBytes - 16
			for index := 0; index < lengthEncryptedMessage/aes.BlockSize; index++ {
				chunk := make([]byte, aes.BlockSize)
				decryptCBCBlock.CryptBlocks(chunk, tempBuffer[16:][aes.BlockSize*index:(index+1)*aes.BlockSize])
				decryptedBufferPadded = append(decryptedBufferPadded, chunk...)
			}
			lengthDecryptedBufferPadded := len(decryptedBufferPadded)
			if lengthDecryptedBufferPadded > 0 && lengthDecryptedBufferPadded <= lengthEncryptedMessage {
				unPaddedBuffer := PKCS7UnPadding(decryptedBufferPadded[:lengthEncryptedMessage])
				copy(buffer, unPaddedBuffer)
				return len(unPaddedBuffer), connectionError
			}
		}
	}
	return numberOfReceivedBytes, connectionError
}

type TunnelWriter struct {
	CipherBlock        cipher.Block
	ActiveSocketWriter SocketWriter
}

func (tunnelWriter *TunnelWriter) Write(buffer []byte) (int, error) {
	paddedBuffer := PKCS7PaddingAES(buffer)
	encryptedBuffer := make([]byte, 0)
	// Always create a new IV for each message
	iv, _ := GenerateIV()
	encryptCBCBlock := cipher.NewCBCEncrypter(tunnelWriter.CipherBlock, iv)
	for index := 0; index < len(paddedBuffer)/aes.BlockSize; index++ {
		chunk := make([]byte, aes.BlockSize)
		encryptCBCBlock.CryptBlocks(chunk, paddedBuffer[aes.BlockSize*index:(index+1)*aes.BlockSize])
		encryptedBuffer = append(encryptedBuffer, chunk...)
	}
	return tunnelWriter.ActiveSocketWriter.Write(append(iv, encryptedBuffer...))
}

func (tunnelWriter *TunnelWriter) Flush() error {
	return tunnelWriter.ActiveSocketWriter.Flush()
}
