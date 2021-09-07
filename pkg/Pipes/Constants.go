package Pipes

import (
	"io"
	"net"
)

func closer(conn1, conn2 io.Closer) {
	_ = conn1.Close()
	_ = conn2.Close()
}

func ForwardTraffic(clientConnection net.Conn, targetConnection net.Conn) error {
	defer closer(clientConnection, targetConnection)
	go io.Copy(clientConnection, targetConnection)
	_, forwardError := io.Copy(targetConnection, clientConnection)
	return forwardError
}
