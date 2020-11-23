package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"log"
)

func SetupLocalForward(host *string, port *string, masterHost *string, masterPort *string) {
	proxy := new(PortForward.LocalForward)
	proxy.TargetHost = *host
	proxy.TargetPort = *port
	proxy.SetLoggingMethod(log.Print)
	ControllersSetup.GeneralSlave(masterHost, masterPort, proxy)
}
