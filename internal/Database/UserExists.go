package Database

import (
	"database/sql"
	"log"
)

func Exists(connection *sql.DB, username *string) int {
	result, queryingError := connection.Query("SELECT id FROM users WHERE username = ?", *username)
	if queryingError != nil {
		log.Fatal(queryingError)
	}
	defer result.Close()
	if result.Next() {
		userId := new(int)
		scanError := result.Scan(userId)
		if scanError != nil {
			log.Fatal(scanError)
		}
		return *userId
	}
	return -1
}
