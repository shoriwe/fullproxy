package network

import "net"

func ListenAny() net.Listener {
	l, _ := net.Listen("tcp", "localhost:0")
	return l
}
