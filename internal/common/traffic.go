package common

import (
	"io"
	"net"
)

const (
	SniffSeparator = "\n\n------------------\n\n"
)

func closer(conn1, conn2 io.Closer) {
	_ = conn1.Close()
	_ = conn2.Close()
}

func ForwardTraffic(
	clientConnection net.Conn, targetConnection io.ReadWriteCloser) error {
	defer closer(clientConnection, targetConnection)
	go io.Copy(clientConnection, targetConnection)
	_, err := io.Copy(targetConnection, clientConnection)
	return err
}
