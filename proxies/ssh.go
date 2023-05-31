package proxies

import (
	"io"
	"net"

	"github.com/shoriwe/fullproxy/v3/sshd"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Network  string
	Address  string
	Listener net.Listener
	Dial     network.DialFunc
	Client   *ssh.Client
}

func (s *SSH) Handle(client net.Conn) {
	defer client.Close()
	target, dErr := s.Dial(s.Network, s.Address)
	if dErr == nil {
		defer target.Close()
		go io.Copy(client, target)
		io.Copy(target, client)
	}
}

func (s *SSH) Close() {
	s.Listener.Close()
	s.Client.Close()
}

func (s *SSH) Addr() net.Addr {
	return s.Listener.Addr()
}

func (s *SSH) Serve() error {
	go sshd.KeepAlive(s.Client)
	for {
		client, err := s.Listener.Accept()
		if err == nil {
			go s.Handle(client)
		}
	}
}
