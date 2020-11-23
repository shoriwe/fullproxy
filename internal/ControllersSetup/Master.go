package ControllersSetup

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers/Master"
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
	controller := new(Master.RemotePortForward)

	controller.TLSConfiguration = tlsConfiguration
	controller.MasterConnection = masterConnection
	controller.MasterConnectionReader = masterConnectionReader
	controller.MasterConnectionWriter = masterConnectionWriter
	controller.Server = server
	controller.RemoteHost = *remoteHost
	controller.RemotePort = *remotePort
	controller.Tries = tries
	controller.Timeout = timeout
	log.Fatal(controller.Serve())
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

	controller := new(Master.General)
	controller.Server = server
	controller.MasterConnection = masterConnection
	controller.MasterConnectionReader = masterConnectionReader
	controller.MasterConnectionWriter = masterConnectionWriter
	controller.TLSConfiguration = tlsConfiguration
	controller.MasterHost = *host
	controller.Tries = tries
	controller.Timeout = timeout
	controller.SetLoggingMethod(log.Print)
	log.Fatal(controller.Serve())
}
