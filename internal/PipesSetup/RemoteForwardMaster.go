package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Master"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func RemoteForwardMaster(
	masterHost *string, masterPort *string,
	remoteHost *string, remotePort *string,
	tries *int, timeout *time.Duration) {
	server, bindError := Sockets.Bind(masterHost, masterPort)
	if bindError != nil {
		log.Fatal(bindError)
	}
	log.Print("Bind successfully in: ", *masterHost, ":", *masterPort)
	tlsConfiguration, creationError := Sockets.CreateMasterTLSConfiguration()
	if creationError != nil {
		log.Fatal(creationError)
	}
	masterConnection, connectionError := server.Accept()
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	log.Println("Slave connected from: ", masterConnection.RemoteAddr().String())
	masterConnection = Sockets.UpgradeServerToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	pipe := new(Master.RemotePortForward)

	pipe.TLSConfiguration = tlsConfiguration
	pipe.MasterConnection = masterConnection
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.Server = server
	pipe.SlaveHost = *remoteHost
	pipe.SlavePort = *remotePort
	_ = pipe.SetTries(*tries)
	_ = pipe.SetTimeout(*timeout)
	Pipes.Serve(pipe)
}
