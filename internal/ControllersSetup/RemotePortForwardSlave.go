package ControllersSetup

import (
	"github.com/shoriwe/FullProxy/pkg/ConnectionControllers/Slave"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func RemotePortForwardSlave(
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
	controller := new(Slave.RemotePortForward)
	controller.MasterHost = *socks5Host
	controller.MasterPort = *socks5Port
	controller.TLSConfiguration = tlsConfiguration
	controller.LocalServer = localServer
	controller.MasterConnection = masterConnection
	controller.MasterConnectionReader = masterConnectionReader
	controller.MasterConnectionWriter = masterConnectionWriter
	controller.Tries = tries
	controller.Timeout = timeout
	controller.SetLoggingMethod(log.Print)
	log.Fatal(controller.Serve())
}
