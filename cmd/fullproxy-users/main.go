package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	Templates "github.com/shoriwe/FullProxy/internal/Tools"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		_, _ = fmt.Fprintf(os.Stderr, "%s COMMAND DATABASE_FILE USERNAME\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stderr, "Available commands:\n"+
			"\t- new\n"+
			"\t- delete\n"+
			"\t- set")
		os.Exit(0)
	}
	stdinReader := bufio.NewReader(os.Stdin)
	switch os.Args[1] {
	case "new":
		fmt.Print("Password: ")
		password, _ := stdinReader.ReadString('\n')
		if os.PathSeparator == '\\' {
			password = password[:len(password)-1]
		}
		password = password[:len(password)-1]
		passwordHash := Templates.SHA3512([]byte(password))
		file, creationError := os.Create(os.Args[2])
		if creationError != nil {
			panic(creationError)
		}
		users := map[string]string{
			os.Args[3]: passwordHash,
		}
		bytes, marshalError := json.Marshal(users)
		if marshalError != nil {
			panic(marshalError)
		}
		_, _ = file.Write(bytes)
		_ = file.Close()
	case "set":
		fmt.Print("Password: ")
		password, _ := stdinReader.ReadString('\n')

		if os.PathSeparator == '\\' {
			password = password[:len(password)-1]
		}
		password = password[:len(password)-1]

		passwordHash := Templates.SHA3512([]byte(password))
		users := Templates.LoadUsers(os.Args[2])
		users[os.Args[3]] = passwordHash
		bytes, marshalError := json.Marshal(users)
		if marshalError != nil {
			panic(marshalError)
		}
		file, openError := os.Create(os.Args[2])
		if openError != nil {
			panic(openError)
		}
		_, _ = file.Write(bytes)
		_ = file.Close()
	case "delete":
		users := Templates.LoadUsers(os.Args[2])
		_, ok := users[os.Args[3]]
		if !ok {
			panic("User not found")
		}
		delete(users, os.Args[3])
		bytes, _ := json.Marshal(users)
		file, openError := os.Create(os.Args[2])
		if openError != nil {
			panic(openError)
		}
		_, _ = file.Write(bytes)
		_ = file.Close()
	}
}
