package CryptoTools

import (
	"crypto/sha512"
)


func Sha3_512_256(buffer []byte) []byte{
	return sha512.New512_256().Sum(buffer)
}


func GetPasswordHashPasswordByteArray(username *[]byte, password *[]byte) []byte{
	if *username != nil && password  != nil{
		return Sha3_512_256(*password)
	}
	return []byte{}
}


func GetPasswordHashPasswordString(username *string, password *string) []byte{
	if len(*username) > 0{
		return Sha3_512_256([]byte(*password))
	}
	return []byte{}
}