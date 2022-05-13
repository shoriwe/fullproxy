package main

import (
	"crypto/tls"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"log"
	"os"
)

func slave() {
	if len(os.Args) == 2 {
		printAndExit("Usage: fullproxy slave MASTER_ADDRESS", 0)
	}
	listen := os.Args[2]
	if listen == "" {
		printAndExit("no master address provided", 1)
	}
	slaveListener, newSlaveError := listeners.NewSlave("tcp", listen, &tls.Config{InsecureSkipVerify: true})
	if newSlaveError != nil {
		printAndExit(newSlaveError.Error(), 1)
	}
	log.Fatal(slaveListener.Serve())
}
