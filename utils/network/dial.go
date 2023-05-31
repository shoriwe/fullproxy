package network

import "net"

type CloseFunc func() error

func NopClose() error {
	return nil
}

type DialFunc func(n, a string) (net.Conn, error)

func Dial(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}
