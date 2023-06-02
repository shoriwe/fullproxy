package network

import "net"

type ListenFunc func(n, addr string) (net.Listener, error)

func ListenAny() net.Listener {
	l, _ := net.Listen("tcp", "localhost:0")
	return l
}
