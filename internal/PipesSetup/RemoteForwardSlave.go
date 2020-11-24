package PipesSetup

import (
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Slave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func RemoteForwardSlave(
	socks5Host *string, socks5Port *string,
	bindHost *string, bindPort *string,
	tries *int, timeout *time.Duration,
	inboundLists [2]string) {
	tlsConfiguration, configurationError := Sockets.CreateSlaveTLSConfiguration()
	if configurationError != nil {
		log.Fatal(configurationError)
	}
	masterConnection, connectionError := Sockets.Connect(socks5Host, socks5Port)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	log.Println("Connected to master: ", masterConnection.RemoteAddr().String())
	localServer, bindError := Sockets.Bind(bindHost, bindPort)
	if bindError != nil {
		log.Fatal(bindError)
	}
	log.Print("Bind successfully in: ", *bindHost, ":", *bindPort)
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
	_ = pipe.SetTries(*tries)
	_ = pipe.SetTimeout(*timeout)
	_ = pipe.SetLoggingMethod(log.Print)
	filter, loadingError := IOTools.LoadList(inboundLists[0], inboundLists[1])
	if loadingError == nil {
		_ = pipe.SetInboundFilter(filter)
	}
	Pipes.Serve(pipe)
}
