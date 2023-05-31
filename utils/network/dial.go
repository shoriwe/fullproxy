package network

import "net"

type DialFunc func(n, a string) (net.Conn, error)

func Dial(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}
