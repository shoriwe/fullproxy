package Templates

import (
	"golang.org/x/crypto/sha3"
)

func SHA3512(buffer []byte) []byte {
	hashHandler := sha3.New512()
	hashHandler.Write(buffer)
	return hashHandler.Sum(nil)
}
