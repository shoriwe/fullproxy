package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Master"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func MasterRemote(
	masterHost *string, masterPort *string,
	remoteHost *string, remotePort *string,
	tries int, timeout time.Duration) {
	server, bindError := Sockets.Bind(masterHost, masterPort)
	if bindError != nil {
		log.Fatal(bindError)
	}
	tlsConfiguration, creationError := Sockets.CreateMasterTLSConfiguration()
	if creationError != nil {
		log.Fatal(creationError)
	}
	masterConnection, connectionError := server.Accept()
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	masterConnection = Sockets.UpgradeServerToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	pipe := new(Master.RemotePortForward)

	pipe.TLSConfiguration = tlsConfiguration
	pipe.MasterConnection = masterConnection
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.Server = server
	pipe.RemoteHost = *remoteHost
	pipe.RemotePort = *remotePort
	pipe.Tries = tries
	pipe.Timeout = timeout
	Pipes.Serve(pipe)
}

func MasterGeneral(host *string, port *string, tries int, timeout time.Duration) {
	server, bindError := Sockets.Bind(host, port)
	if bindError != nil {
		log.Fatal(bindError)
	}
	tlsConfiguration, creationError := Sockets.CreateMasterTLSConfiguration()
	if creationError != nil {
		log.Fatal(creationError)
	}
	masterConnection, connectionError := server.Accept()
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	masterConnection = Sockets.UpgradeServerToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)

	pipe := new(Master.General)
	pipe.Server = server
	pipe.MasterConnection = masterConnection
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.TLSConfiguration = tlsConfiguration
	pipe.MasterHost = *host
	pipe.Tries = tries
	pipe.Timeout = timeout
	pipe.SetLoggingMethod(log.Print)
	Pipes.Serve(pipe)
}
