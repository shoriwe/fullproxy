package Database

import (
	"database/sql"
	"encoding/hex"
	"github.com/shoriwe/FullProxy/internal/Templates"
	"log"
)

func Update(databaseFilePath string, username string, newPassword string) {
	connection, connectionError := sql.Open("sqlite3", databaseFilePath)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	defer connection.Close()
	userId := Exists(connection, &username)
	if userId == -1 {
		log.Fatal("User doesn't Exists in the Database")
	}
	newPasswordHash := hex.EncodeToString(Templates.SHA3512([]byte(newPassword)))
	_, updateError := connection.Exec("UPDATE users SET password = ? WHERE id = ?", newPasswordHash, userId)
	if updateError != nil {
		log.Fatal(updateError)
	}
	log.Print("User updated successfully")
}
