package proxies

import (
	"io"
	"net"
)

type Forward struct {
	Network  string
	Address  string
	Listener net.Listener
	Dial     func(network, address string) (net.Conn, error)
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

func (f *Forward) Addr() net.Addr {
	return f.Listener.Addr()
}

func (f *Forward) Close() {
	f.Listener.Close()
}

func (f *Forward) Serve() (err error) {
	var client net.Conn
	for client, err = f.Listener.Accept(); err == nil; client, err = f.Listener.Accept() {
		go f.Handle(client)
	}
	return err
}
