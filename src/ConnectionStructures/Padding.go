package ConnectionStructures

import (
	"bytes"
	"crypto/aes"
)

func PKCS7PaddingAES(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length > 0 {
		unPadding := int(origData[length-1])
		if length > unPadding {
			return origData[:(length - unPadding)]
		}
	}
	return origData
}
