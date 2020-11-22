package PortForward

import (
	"bufio"
	"fmt"
	"github.com/shoriwe/FullProxy/pkg/ConnectionStructures"
	"github.com/shoriwe/FullProxy/pkg/MasterSlave"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

func CreateLocalPortForwardSession(
	clientConnection net.Conn,
	clientReader *bufio.Reader,
	clientWriter *bufio.Writer,
	args ...interface{}) {
	targetAddress := args[0].(*string)
	targetPort := args[1].(*string)
	targetConnection := Sockets.Connect(targetAddress, targetPort)
	if targetConnection != nil {
		targetReader, targetWriter := ConnectionStructures.CreateSocketConnectionReaderWriter(targetConnection)
		Basic.Proxy(clientConnection, targetConnection, clientReader, clientWriter, targetReader, targetWriter)
	} else {
		_ = clientConnection.Close()
	}
}

func StartLocalPortForward(targetAddress *string, targetPort *string, masterAddress *string, masterPort *string) {
	if !(*targetAddress == "" || *targetPort == "" || *masterAddress == "" || *masterPort == "") {
		MasterSlave.GeneralSlave(masterAddress, masterPort, CreateLocalPortForwardSession, targetAddress, targetPort)
	} else {
		fmt.Println("All flags need to be in use")
	}
}
