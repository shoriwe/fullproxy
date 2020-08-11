package PortForward


import (
	"fmt"
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/MasterSlave"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
)


func CreatePortForwardSession(masterConnection net.Conn, masterReader ConnectionStructures.SocketReader, masterWriter ConnectionStructures.SocketWriter, args...interface{}){
	targetAddress := args[0].(*string)
	targetPort := args[1].(*string)
	targetConnection := Sockets.Connect(targetAddress, targetPort)
	if targetConnection != nil{
		targetReader, targetWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(targetConnection)
		Basic.Proxy(masterConnection, targetConnection, masterReader, masterWriter, targetReader, targetWriter)
	} else {
		_ = masterConnection.Close()
	}
}


func StartPortForward(targetAddress *string, targetPort *string, masterAddress *string, masterPort *string){
	if !(*targetAddress == "" || *targetPort == "" || *masterAddress == "" || *masterPort == ""){
		MasterSlave.Slave(masterAddress, masterPort, CreatePortForwardSession, targetAddress, targetPort)
	} else {
		fmt.Println("All flags need to be in use")
	}
}
