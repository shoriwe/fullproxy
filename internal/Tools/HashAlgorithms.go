package Templates

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

func SHA3512(buffer []byte) string {
	hashHandler := sha3.New512()
	hashHandler.Write(buffer)
	return hex.EncodeToString(hashHandler.Sum(nil))
}
