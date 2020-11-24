package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
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
	PipesSetup.GeneralSlave(masterHost, masterPort, proxy)
}
