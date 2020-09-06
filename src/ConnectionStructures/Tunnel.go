package ConnectionStructures

import (
	"crypto/aes"
	"crypto/cipher"
)


type TunnelReader struct {
	CipherBlock cipher.Block
	ActiveSocketReader SocketReader
}


func (tunnelReader *TunnelReader) Read(buffer []byte) (int,  error){
	tempBuffer := make([]byte, len(buffer))
	numberOfReceivedBytes, connectionError := tunnelReader.ActiveSocketReader.Read(tempBuffer)
	if connectionError == nil{
		if numberOfReceivedBytes > 0 {
			decryptedBufferPadded := make([]byte, 0)
			for index := 0; index < numberOfReceivedBytes/aes.BlockSize; index++ {
				chunk := make([]byte, aes.BlockSize)
				tunnelReader.CipherBlock.Decrypt(chunk, tempBuffer[aes.BlockSize*index:(index+1)*aes.BlockSize])
				decryptedBufferPadded = append(decryptedBufferPadded, chunk...)
			}
			if lengthDecryptedBufferPadded := len(decryptedBufferPadded); lengthDecryptedBufferPadded > 0 && lengthDecryptedBufferPadded <= numberOfReceivedBytes{
				unPaddedBuffer := PKCS7UnPadding(decryptedBufferPadded[:numberOfReceivedBytes])
				copy(buffer, unPaddedBuffer)
				return len(unPaddedBuffer), connectionError
			}
		}
	}
	return numberOfReceivedBytes, connectionError
}


type TunnelWriter struct {
	CipherBlock cipher.Block
	ActiveSocketWriter SocketWriter
}


func (tunnelWriter *TunnelWriter) Write(buffer []byte) (int, error){
	paddedBuffer := PKCS7PaddingAES(buffer)
	encryptedBuffer := make([]byte, 0)
	for index := 0; index < len(paddedBuffer) / aes.BlockSize; index ++{
		chunk := make([]byte, aes.BlockSize)
		tunnelWriter.CipherBlock.Encrypt(chunk, paddedBuffer[aes.BlockSize * index:(index+1)*aes.BlockSize])
		encryptedBuffer = append(encryptedBuffer, chunk...)
	}
	return tunnelWriter.ActiveSocketWriter.Write(encryptedBuffer)
}

func (tunnelWriter *TunnelWriter) Flush() error {
	return tunnelWriter.ActiveSocketWriter.Flush()
}


