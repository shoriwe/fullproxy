package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/SetupControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
)

func SetupRemoteForward(host *string, port *string, masterHost *string, masterPort *string) {
	tlsConfiguration, configurationError := Sockets.CreateSlaveTLSConfiguration()
	if configurationError != nil {
		log.Fatal(configurationError)
	}
	proxy := new(PortForward.RemoteForward)
	proxy.TLSConfiguration = tlsConfiguration
	proxy.MasterHost = *masterHost
	proxy.MasterPort = *masterPort
	proxy.SetLoggingMethod(log.Print)
	SetupControllers.RemotePortForwardSlave(masterHost, masterPort, host, port)
}
