package Tools

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"io"
	"os"
)

func LoadUsers(filePath string) map[string]string {
	file, openError := os.Open(filePath)
	if openError != nil {
		panic(openError)
	}
	defer file.Close()
	content, readError := io.ReadAll(file)
	if readError != nil {
		panic(readError)
	}
	result := map[string]string{}
	marshalError := json.Unmarshal(content, &result)
	if marshalError != nil {
		panic(marshalError)
	}
	return result
}

func SHA3512(buffer []byte) string {
	hashHandler := sha3.New512()
	hashHandler.Write(buffer)
	return hex.EncodeToString(hashHandler.Sum(nil))
}
