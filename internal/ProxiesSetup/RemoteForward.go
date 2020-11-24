package ProxiesSetup

import (
	"github.com/shoriwe/FullProxy/internal/PipesSetup"
	"time"
)

func SetupRemoteForward(
	host *string, port *string,
	masterHost *string, masterPort *string,
	tries *int, timeout *time.Duration,
	inboundLists [2]string) {
	PipesSetup.RemoteForwardSlave(masterHost, masterPort, host, port, tries, timeout, inboundLists)
}
