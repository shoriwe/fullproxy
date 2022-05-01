package common

import (
	"io"
	"net"
)

func closer(conn1, conn2 io.Closer) {
	_ = conn1.Close()
	_ = conn2.Close()
}

func netCopy(dst, src net.Conn) error {
	defer src.Close()
	defer dst.Close()
	_, err := io.Copy(dst, src)
	return err
}

func ForwardTraffic(clientConnection net.Conn, targetConnection net.Conn) error {
	defer closer(clientConnection, targetConnection)
	go netCopy(clientConnection, targetConnection)
	return netCopy(targetConnection, clientConnection)
}
