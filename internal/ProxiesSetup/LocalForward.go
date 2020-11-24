package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/IOTools"
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"github.com/shoriwe/FullProxy/pkg/Proxies/PortForward"
	"log"
	"time"
)

func SetupLocalForward(
	host *string, port *string,
	masterHost *string, masterPort *string,
	tries *int, timeout *time.Duration,
	inboundLists [2]string) {
	proxy := new(PortForward.LocalForward)
	proxy.TargetHost = *host
	proxy.TargetPort = *port
	_ = proxy.SetTries(*tries)
	_ = proxy.SetTimeout(*timeout)
	_ = proxy.SetLoggingMethod(log.Print)
	filter, loadingError := IOTools.LoadList(inboundLists[0], inboundLists[1])
	if loadingError == nil {
		_ = proxy.SetInboundFilter(filter)
	}
	PipesSetup.GeneralSlave(masterHost, masterPort, proxy)
}
