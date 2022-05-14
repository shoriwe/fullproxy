package main

import (
	"crypto/tls"
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"log"
	"net/url"
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
	listenURL, parseError := url.Parse(listen)
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	listenAddress := fmt.Sprintf("%s:%s", listenURL.Hostname(), listenURL.Port())
	slaveListener, newSlaveError := listeners.NewSlave(listenURL.Scheme, listenAddress, &tls.Config{InsecureSkipVerify: true})
	if newSlaveError != nil {
		printAndExit(newSlaveError.Error(), 1)
	}
	log.Fatal(slaveListener.Serve())
}
