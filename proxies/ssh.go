package proxies

import (
	"io"
	"net"
	"time"

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

func (s *SSH) keepAlive() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, _, err := s.Client.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				return
			}
		}
	}
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
	go s.keepAlive()
	for {
		client, err := s.Listener.Accept()
		if err == nil {
			go s.Handle(client)
		}
	}
}
