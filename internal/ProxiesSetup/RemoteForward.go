package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"log"
	"time"
)

func SetupRemoteForward(
	host *string, port *string,
	masterHost *string, masterPort *string,
	tries int, timeout time.Duration) {
	tlsConfiguration, configurationError := Sockets.CreateSlaveTLSConfiguration()
	if configurationError != nil {
		log.Fatal(configurationError)
	}
	proxy := new(PortForward.RemoteForward)
	proxy.TLSConfiguration = tlsConfiguration
	proxy.MasterHost = *masterHost
	proxy.MasterPort = *masterPort
	proxy.SetTries(tries)
	proxy.SetTimeout(timeout)
	proxy.SetLoggingMethod(log.Print)
	ControllersSetup.RemotePortForwardSlave(masterHost, masterPort, host, port, tries, timeout)
}
