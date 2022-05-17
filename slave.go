package main

import (
	"crypto/tls"
	"flag"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"log"
	"net/url"
	"os"
)

func slave() {
	var (
		master string
	)
	translateCmd := flag.NewFlagSet("translate", flag.ExitOnError)
	translateCmd.StringVar(&master, "master", "", "Address of master server. Argument URL structure is 'network://host:port'")
	cmdParseError := translateCmd.Parse(os.Args[2:])
	if cmdParseError != nil {
		printAndExit(cmdParseError.Error(), 1)
	}
	masterURL, parseError := url.Parse(master)
	if parseError != nil {
		printAndExit(parseError.Error(), 1)
	}
	slaveListener, newSlaveError := listeners.NewSlave(
		masterURL.Scheme,
		masterURL.Host,
		&tls.Config{InsecureSkipVerify: true},
	)
	if newSlaveError != nil {
		printAndExit(newSlaveError.Error(), 1)
	}
	log.Fatal(slaveListener.Serve())
}
