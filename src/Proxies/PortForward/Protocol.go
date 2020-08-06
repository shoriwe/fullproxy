package PortForward

import (
	"github.com/shoriwe/FullProxy/src/Interface"
	"github.com/shoriwe/FullProxy/src/Proxies/Basic"
	"github.com/shoriwe/FullProxy/src/Sockets"
	"net"
)

func CreatePortForwardSession(masterConnection net.Conn, args...interface{}){
	targetAddress := args[0].(*string)
	targetPort := args[1].(*string)
	targetConnection := Sockets.Connect(targetAddress, targetPort)
	if targetConnection != nil{
		masterReader, masterWriter := Sockets.CreateReaderWriter(masterConnection)
		targetReader, targetWriter := Sockets.CreateReaderWriter(targetConnection)
		Basic.Proxy(masterConnection, targetConnection, masterReader, masterWriter, targetReader, targetWriter)
	} else {
		_ = masterConnection.Close()
	}
}


func  StartPortForward(targetAddress *string, targetPort *string, masterAddress *string, masterPort *string){
	if *targetAddress == "" || *targetPort == "" || *masterAddress == "" || *masterPort == ""{

	}
	Interface.Slave(targetAddress, masterPort, CreatePortForwardSession, targetAddress, targetPort)
}
