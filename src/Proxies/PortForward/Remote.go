package PortForward

import (
	"fmt"
	"github.com/shoriwe/FullProxy/src/MasterSlave"
)

func StartRemotePortForward(localAddress *string, localPort *string, masterAddress *string, masterPort *string) {
	if !(*localAddress == "" || *localPort == "" || *masterAddress == "" || *masterPort == "") {
		MasterSlave.RemotePortForwardSlave(masterAddress, masterPort, localAddress, localPort)
	} else {
		fmt.Println("All flags need to be in use")
	}
}
