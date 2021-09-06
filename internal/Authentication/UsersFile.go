package Authentication

import (
	"encoding/json"
	"github.com/shoriwe/FullProxy/internal/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"io"
	"os"
)

func UsersFile(usersFiles string) Types.AuthenticationMethod {
	file, openError := os.Open(usersFiles)
	if openError != nil {
		panic(openError)
	}
	content, readingError := io.ReadAll(file)
	if readingError != nil {
		panic(readingError)
	}
	// Local credentials
	var entries map[string]string
	unMarshalError := json.Unmarshal(content, &entries)
	if unMarshalError != nil {
		panic(unMarshalError)
	}
	return func(username []byte, password []byte) (bool, error) {
		dbPassword, found := entries[string(username)]
		if !found {
			return false, nil
		}
		if dbPassword == Templates.SHA3512(password) {
			return false, nil
		}
		return true, nil
	}
}
