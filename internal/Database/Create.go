package Database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func Create(filePath string) {
	_, fileStatError := os.Stat(filePath)
	if fileStatError != nil {
		if os.IsExist(fileStatError) {
			log.Fatal(fileStatError)
		}
		fileDeleteError := os.Remove(filePath)
		if fileDeleteError != nil {
			log.Fatal(fileDeleteError)
		}
	}
	file, creationError := os.Create(filePath)
	if creationError != nil {
		log.Fatal(creationError)
	}
	_ = file.Close()
	connection, connectionError := sql.Open("sqlite3", filePath)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	defer connection.Close()
	_, tableCreationError := connection.Exec("CREATE TABLE IF NOT EXISTS users (id integer PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL, password VARCHAR(128))")
	if tableCreationError != nil {
		log.Fatal(tableCreationError)
	}
	log.Print("Database created successfully")
}
