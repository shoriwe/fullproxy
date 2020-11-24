package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Slave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func RemoteForwardSlave(
	socks5Host *string, socks5Port *string,
	bindHost *string, bindPort *string,
	tries int, timeout time.Duration) {
	tlsConfiguration, configurationError := Sockets.CreateSlaveTLSConfiguration()
	if configurationError != nil {
		log.Fatal(configurationError)
	}
	masterConnection, connectionError := Sockets.Connect(socks5Host, socks5Port)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	localServer, bindError := Sockets.Bind(bindHost, bindPort)
	if bindError != nil {
		log.Fatal(bindError)
	}
	masterConnection = Sockets.UpgradeClientToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	pipe := new(Slave.RemotePortForward)
	pipe.MasterHost = *socks5Host
	pipe.MasterPort = *socks5Port
	pipe.TLSConfiguration = tlsConfiguration
	pipe.LocalServer = localServer
	pipe.MasterConnection = masterConnection
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.Tries = tries
	pipe.Timeout = timeout
	Pipes.Serve(pipe)
}
