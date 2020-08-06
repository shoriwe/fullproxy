package ConnectionStructures

import (
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
		decryptedBufferPadded := make([]byte, len(tempBuffer))
		tunnelReader.CipherBlock.Decrypt(decryptedBufferPadded, tempBuffer)
		unPaddedBuffer := PKCS7UnPadding(decryptedBufferPadded[:numberOfReceivedBytes])
		copy(buffer, decryptedBufferPadded[:numberOfReceivedBytes])
		return len(unPaddedBuffer), connectionError
	}
	return numberOfReceivedBytes, connectionError
}


type TunnelWriter struct {
	CipherBlock cipher.Block
	ActiveSocketWriter SocketWriter
}


func (tunnelWriter *TunnelWriter) Write(buffer []byte) (int, error){
	paddedBuffer := PKCS7PaddingAES(buffer)
	encryptedBuffer := make([]byte, len(paddedBuffer))
	tunnelWriter.CipherBlock.Encrypt(encryptedBuffer, paddedBuffer)
	return tunnelWriter.ActiveSocketWriter.Write(encryptedBuffer)
}

func (tunnelWriter *TunnelWriter) Flush() error {
	return tunnelWriter.ActiveSocketWriter.Flush()
}


