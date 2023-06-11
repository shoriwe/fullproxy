package circuit

import (
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

func TestConn_Close(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	go listener.Accept()
	conn := network.Dial(listener.Addr().String())
	defer conn.Close()
	c := Conn{
		Conn:           conn,
		CloseFunctions: []network.CloseFunc{network.NopClose},
	}
	defer c.Close()
}
