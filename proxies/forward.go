package proxies

import (
	"io"
	"net"
)

type Forward struct {
	Network string
	Address string
	Accept  func() (net.Conn, error)
	Dial    func(network, address string) (net.Conn, error)
}

func (f *Forward) Handle(client net.Conn) {
	defer client.Close()
	target, dErr := f.Dial(f.Network, f.Address)
	if dErr == nil {
		defer target.Close()
		go io.Copy(client, target)
		io.Copy(target, client)
	}
}

func (f *Forward) Close() {}

func (f *Forward) Serve() {
	for {
		client, aErr := f.Accept()
		if aErr == nil {
			go f.Handle(client)
		}
	}
}
