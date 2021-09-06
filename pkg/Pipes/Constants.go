package Pipes

import (
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"io"
	"log"
	"net"
)

var (
	NewConnection = []byte{1}
	// Shutdown              = []byte{2}
	FailToConnectToTarget = []byte{3}
	UnknownOperation      = []byte{4}
)

func Serve(pipe Types.Pipe) {
	log.Fatal(pipe.Serve())
}

func closer(conn1, conn2 io.Closer) {
	conn1.Close()
	conn2.Close()
}

func ForwardTraffic(clientConnection net.Conn, targetConnection net.Conn) error {
	defer closer(clientConnection, targetConnection)
	go io.Copy(clientConnection, targetConnection)
	_, writeError := io.Copy(targetConnection, clientConnection)
	return writeError
}
