package PipesSetup

import (
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Slave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"log"
)

func GeneralSlave(host *string, port *string, proxy Types.ProxyProtocol) {
	masterConnection, connectionError := Sockets.Connect(host, port)
	if connectionError != nil {
		log.Fatal(connectionError)
	}
	log.Println("Connected to master: ", masterConnection.RemoteAddr().String())
	tlsConfiguration, creationError := Sockets.CreateSlaveTLSConfiguration()
	if creationError != nil {
		log.Fatal(creationError)
	}
	masterConnection = Sockets.UpgradeClientToTLS(masterConnection, tlsConfiguration)
	masterConnectionReader, masterConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(masterConnection)
	pipe := new(Slave.General)
	pipe.MasterConnection = masterConnection
	pipe.TLSConfiguration = tlsConfiguration
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.MasterHost = *host
	pipe.MasterPort = *port
	pipe.ProxyProtocol = proxy
	_ = pipe.SetLoggingMethod(log.Print)
	Pipes.Serve(pipe)
}
