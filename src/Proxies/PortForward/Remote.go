package PortForward

import "github.com/shoriwe/FullProxy/src/MasterSlave"

func StartRemotePortForward(localAddress *string, localPort *string, masterAddress *string, masterPort *string){
	MasterSlave.RemotePortForwardSlave(masterAddress, masterPort, localAddress, localPort)
}
