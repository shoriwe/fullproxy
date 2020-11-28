package Database

import (
	"database/sql"
	"encoding/hex"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"log"
)

func Add(databaseFilePath string, username string, password string) {
	connection, connectionError := sql.Open("sqlite3", databaseFilePath)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	defer connection.Close()
	if Exists(connection, &username) != -1 {
		log.Fatal("User already exists in the database")
	}
	passwordHash := hex.EncodeToString(Templates.SHA3512([]byte(password)))
	_, insertError := connection.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, passwordHash)
	if insertError != nil {
		log.Fatal(insertError)
	}
	log.Print("User added successfully")
}
