package Authentication

import (
	"database/sql"
	"encoding/hex"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"log"
)

func SQLite3Authentication(databaseFilePath string) Types.AuthenticationMethod {
	connection, connectionError := sql.Open("sqlite3", databaseFilePath)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	return func(username []byte, password []byte) bool {
		passwordHash := hex.EncodeToString(Templates.SHA3512(password))
		result, queryingError := connection.Query("SELECT id FROM users WHERE username = ? AND password = ?", string(username), passwordHash)
		if queryingError != nil {
			log.Print(queryingError)
			return false
		}
		defer result.Close()
		return result.Next()
	}
}
