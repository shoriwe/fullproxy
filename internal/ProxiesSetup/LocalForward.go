package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/SetupControllers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
)

func SetupLocalForward(host *string, port *string, masterHost *string, masterPort *string) {
	proxy := new(PortForward.LocalForward)
	proxy.TargetHost = *host
	proxy.TargetPort = *port
	SetupControllers.GeneralSlave(masterHost, masterPort, proxy)
}
