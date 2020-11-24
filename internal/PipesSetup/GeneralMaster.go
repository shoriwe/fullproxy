package PipesSetup

import (
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Master"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"strings"
	"time"
)

func GeneralMaster(host *string, port *string, tries *int, timeout *time.Duration, inboundLists [2]string) {
	server, bindError := Sockets.Bind(host, port)
	if bindError != nil {
		log.Fatal(bindError)
	}
	log.Print("Bind successfully in: ", *host, ":", *port)
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

	pipe := new(Master.General)
	pipe.Server = server
	pipe.MasterConnection = masterConnection
	pipe.MasterConnectionReader = masterConnectionReader
	pipe.MasterConnectionWriter = masterConnectionWriter
	pipe.TLSConfiguration = tlsConfiguration
	pipe.SlaveHost = strings.Split(masterConnection.RemoteAddr().String(), ":")[0]
	pipe.SetTries(*tries)
	pipe.SetTimeout(*timeout)
	_ = pipe.SetLoggingMethod(log.Print)
	filter, loadingError := IOTools.LoadList(inboundLists[0], inboundLists[1])
	if loadingError == nil {
		_ = pipe.SetInboundFilter(filter)
	}
	Pipes.Serve(pipe)
}
