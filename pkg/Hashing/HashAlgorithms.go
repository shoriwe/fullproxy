package Hashing

import (
	"golang.org/x/crypto/sha3"
)

func PasswordHashingSHA3(buffer []byte) []byte {
	hashHandler := sha3.New512()
	hashHandler.Write(buffer)
	return hashHandler.Sum(nil)
}

func GetPasswordHashPasswordByteArray(username *[]byte, password *[]byte) []byte {
	if len(*username) != 0 && len(*password) != 0 {
		return PasswordHashingSHA3(*password)
	}
	return nil
}
