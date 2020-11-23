package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/ControllersSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"log"
	"time"
)

func SetupLocalForward(
	host *string, port *string,
	masterHost *string, masterPort *string,
	tries int, timeout time.Duration) {
	proxy := new(PortForward.LocalForward)
	proxy.TargetHost = *host
	proxy.TargetPort = *port
	proxy.SetTries(tries)
	proxy.SetTimeout(timeout)
	proxy.SetLoggingMethod(log.Print)
	ControllersSetup.GeneralSlave(masterHost, masterPort, proxy)
}
