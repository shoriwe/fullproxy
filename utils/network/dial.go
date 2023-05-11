package network

import "net"

func Dial(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}
