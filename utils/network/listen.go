package network

import "net"

func ListenAny() net.Listener {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	return l
}
