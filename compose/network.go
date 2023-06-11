package compose

import (
	"fmt"
	"net"

	"github.com/shoriwe/fullproxy/v4/reverse"
	"github.com/shoriwe/fullproxy/v4/sshd"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"golang.org/x/crypto/ssh"
)

const (
	NetworkBasic  = "basic"
	NetworkMaster = "master"
	NetworkSSH    = "ssh"
)

type Network struct {
	Type          string   `yaml:"type" json:"type"`
	Network       *string  `yaml:"network,omitempty" json:"network,omitempty"`
	Address       *string  `yaml:"address,omitempty" json:"address,omitempty"`
	Data          *Network `yaml:"data,omitempty" json:"data,omitempty"`
	Control       *Network `yaml:"control,omitempty" json:"control,omitempty"`
	Auth          *Auth    `yaml:"auth,omitempty" json:"auth,omitempty"`
	Crypto        *Crypto  `yaml:"crypto,omitempty" json:"crypto,omitempty"`
	SlaveListener *bool    `yaml:"slaveListener,omitempty" json:"slaveListener,omitempty"`
	master        *reverse.Master
	sshConn       *ssh.Client
	listener      net.Listener
}

func (n *Network) setupBasicListener(listen network.ListenFunc) (_ net.Listener, err error) {
	if n.listener == nil {
		if n.Network == nil {
			return nil, fmt.Errorf("network not set for basic listener")
		}
		if n.Address == nil {
			return nil, fmt.Errorf("address not set for basic listener")
		}
		n.listener, err = listen(*n.Network, *n.Address)
	}
	return n.listener, err
}

func (n *Network) getMaster() (*reverse.Master, error) {
	if n.master == nil {
		if n.Data == nil {
			return nil, fmt.Errorf("no data listener provided for master")
		}
		if n.Control == nil {
			return nil, fmt.Errorf("no control listener provided for master")
		}
		data, err := n.Data.Listen()
		if err != nil {
			return nil, err
		}
		defer network.CloseOnError(&err, data)
		control, err := n.Control.Listen()
		if err != nil {
			return nil, err
		}
		n.master = &reverse.Master{
			Data:    data,
			Control: control,
		}
	}
	return n.master, nil
}

type slaveListenerWrapper struct {
	*reverse.Master
}

func (slw *slaveListenerWrapper) Accept() (net.Conn, error) {
	return slw.Master.SlaveAccept()
}

func (n *Network) setupMasterListener() (ll net.Listener, err error) {
	master, err := n.getMaster()
	if err != nil {
		return nil, err
	}
	ll = master
	if n.SlaveListener != nil && *n.SlaveListener {
		return &slaveListenerWrapper{master}, nil
	}
	return ll, nil
}

type sshWrapper struct {
	net.Listener
	conn *ssh.Client
}

func (s *sshWrapper) Close() error {
	s.conn.Close()
	return s.Listener.Close()
}

func (n *Network) getSSHConn() (*ssh.Client, error) {
	if n.sshConn == nil {
		if n.Network == nil {
			return nil, fmt.Errorf("network not set for basic listener")
		}
		if n.Address == nil {
			return nil, fmt.Errorf("address not set for basic listener")
		}
		if n.Data == nil {
			return nil, fmt.Errorf("no remote listen configuration")
		}
		if n.Auth == nil {
			return nil, fmt.Errorf("no ssh auth provided")
		}
		var config *ssh.ClientConfig
		config, err := n.Auth.SSHClientConfig()
		if err != nil {
			return nil, err
		}
		n.sshConn, err = ssh.Dial(*n.Network, *n.Address, config)
		if err != nil {
			return nil, err
		}
		go sshd.KeepAlive(n.sshConn)
	}
	return n.sshConn, nil
}

func (n *Network) setupSSHListener() (ll net.Listener, err error) {
	sshConn, err := n.getSSHConn()
	if err == nil {
		defer network.CloseOnError(&err, sshConn)
		data, err := n.Data.setupBasicListener(sshConn.Listen)
		if err == nil {
			ll = &sshWrapper{
				Listener: data,
				conn:     sshConn,
			}
		}
	}
	return ll, err
}

func (n *Network) Listen() (ll net.Listener, err error) {
	switch n.Type {
	case NetworkBasic:
		ll, err = n.setupBasicListener(net.Listen)
	case NetworkMaster:
		ll, err = n.setupMasterListener()
	case NetworkSSH:
		ll, err = n.setupSSHListener()
	default:
		err = fmt.Errorf("unknown network type %s", n.Type)
	}
	if err == nil && n.Crypto != nil {
		ll, err = n.Crypto.WrapListener(ll)
	}
	return ll, err
}

func (n *Network) setupMasterDialFunc() (dialFunc network.DialFunc, err error) {
	master, err := n.getMaster()
	if err == nil {
		dialFunc = master.SlaveDial
	}
	return dialFunc, err
}

func (n *Network) setupSSHDialFunc() (dialFunc network.DialFunc, err error) {
	sshConn, err := n.getSSHConn()
	if err == nil {
		dialFunc = sshConn.Dial
	}
	return dialFunc, err
}

func (n *Network) DialFunc() (network.DialFunc, error) {
	switch n.Type {
	case NetworkBasic:
		return net.Dial, nil
	case NetworkMaster:
		return n.setupMasterDialFunc()
	case NetworkSSH:
		return n.setupSSHDialFunc()
	default:
		return nil, fmt.Errorf("unknown network type %s", n.Type)
	}
}
