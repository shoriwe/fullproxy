package Database

import (
	"database/sql"
	"log"
)

func Delete(databaseFilePath string, username string) {
	connection, connectionError := sql.Open("sqlite3", databaseFilePath)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	defer connection.Close()
	userId := Exists(connection, &username)
	if userId == -1 {
		log.Fatal("User doesn't exists in the database")
	}
	_, deletionError := connection.Exec("DELETE FROM users WHERE id = ?", userId)
	if deletionError != nil {
		log.Fatal(deletionError)
	}
	log.Print("User deleted successfully")
}
