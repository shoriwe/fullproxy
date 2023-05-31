package circuit

import (
	"net"

	"github.com/shoriwe/fullproxy/v3/sshd"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	Network string
	Address string
	Config  ssh.ClientConfig
}

func (s *SSH) Next(dial network.DialFunc) (closeFunc network.CloseFunc, newDial network.DialFunc, err error) {
	var (
		conn       net.Conn
		sshConn    ssh.Conn
		newChannel <-chan ssh.NewChannel
		requests   <-chan *ssh.Request
		client     *ssh.Client
	)
	conn, err = dial(s.Network, s.Address)
	if err == nil {
		defer network.CloseOnError(&err, conn)
		sshConn, newChannel, requests, err = ssh.NewClientConn(conn, "", &s.Config)
		if err == nil {
			go sshd.KeepAlive(client)
			closeFunc = func() error {
				return conn.Close()
			}
			client = ssh.NewClient(sshConn, newChannel, requests)
			newDial = client.Dial
		}
	}
	return closeFunc, newDial, err
}

// newSSH ensures compile time safety, should never be used
func newSSH() Knot {
	return &SSH{}
}
