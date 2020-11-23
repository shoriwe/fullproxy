package ControllersSetup

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers"
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers/Slave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
)

func GeneralSlave(host *string, port *string, proxy ConnectionControllers.ProxyProtocol) {
	masterConnection, connectionError := Sockets.Connect(host, port)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	tlsConfiguration, creationError := Sockets.CreateSlaveTLSConfiguration()
	if creationError != nil {
		log.Fatal(creationError)
	}
	masterConnection = Sockets.UpgradeClientToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	controller := new(Slave.General)
	controller.MasterConnection = masterConnection
	controller.TLSConfiguration = tlsConfiguration
	controller.MasterConnectionReader = masterConnectionReader
	controller.MasterConnectionWriter = masterConnectionWriter
	controller.MasterHost = *host
	controller.MasterPort = *port
	controller.ProxyProtocol = proxy
	controller.SetLoggingMethod(log.Print)
	log.Fatal(controller.Serve())
}
