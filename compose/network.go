package compose

import (
	"fmt"
	"net"

	"github.com/shoriwe/fullproxy/v3/reverse"
	"github.com/shoriwe/fullproxy/v3/sshd"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"golang.org/x/crypto/ssh"
)

const (
	NetworkBasic  = "basic"
	NetworkMaster = "master"
	NetworkSSH    = "ssh"
)

type Network struct {
	Type    string   `yaml:"type" json:"type"`
	Network *string  `yaml:"network,omitempty" json:"network,omitempty"`
	Address *string  `yaml:"address,omitempty" json:"address,omitempty"`
	Data    *Network `yaml:"data,omitempty" json:"data,omitempty"`
	Control *Network `yaml:"control,omitempty" json:"control,omitempty"`
	Auth    *Auth    `yaml:"auth,omitempty" json:"auth,omitempty"`
	Crypto  *Crypto  `yaml:"crypto,omitempty" json:"crypto,omitempty"`
}

func (l *Network) setupBasicListener(listen network.ListenFunc) (net.Listener, error) {
	if l.Network == nil {
		return nil, fmt.Errorf("network not set for basic listener")
	}
	if l.Address == nil {
		return nil, fmt.Errorf("address not set for basic listener")
	}
	return listen(*l.Network, *l.Address)
}

func (l *Network) setupMasterListener() (ll net.Listener, err error) {
	if l.Data == nil {
		return nil, fmt.Errorf("no data listener provided for master")
	}
	if l.Control == nil {
		return nil, fmt.Errorf("no control listener provided for master")
	}
	var (
		data, control net.Listener
	)
	data, err = l.Data.Listen()
	if err != nil {
		return nil, err
	}
	defer network.CloseOnError(&err, data)
	control, err = l.Control.Listen()
	if err != nil {
		return nil, err
	}
	ll = &reverse.Master{
		Data:    data,
		Control: control,
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

func (l *Network) setupSSHListener() (ll net.Listener, err error) {
	if l.Network == nil {
		return nil, fmt.Errorf("network not set for basic listener")
	}
	if l.Address == nil {
		return nil, fmt.Errorf("address not set for basic listener")
	}
	if l.Data == nil {
		return nil, fmt.Errorf("no remote listen configuration")
	}
	if l.Auth == nil {
		return nil, fmt.Errorf("no ssh auth provided")
	}
	var config *ssh.ClientConfig
	config, err = l.Auth.SSHClientConfig()
	if err != nil {
		return nil, err
	}
	var sshConn *ssh.Client
	sshConn, err = ssh.Dial(*l.Network, *l.Address, config)
	if err != nil {
		return nil, err
	}
	go sshd.KeepAlive(sshConn)
	defer network.CloseOnError(&err, sshConn)
	var data net.Listener
	data, err = l.Data.setupBasicListener(sshConn.Listen)
	if err == nil {
		ll = &sshWrapper{
			Listener: data,
			conn:     sshConn,
		}
	}
	return ll, err
}

func (l *Network) Listen() (ll net.Listener, err error) {
	switch l.Type {
	case NetworkBasic:
		ll, err = l.setupBasicListener(net.Listen)
	case NetworkMaster:
		ll, err = l.setupMasterListener()
	case NetworkSSH:
		ll, err = l.setupSSHListener()
	default:
		err = fmt.Errorf("unknown network type %s", l.Type)
	}
	if err == nil && l.Crypto != nil {
		ll, err = l.Crypto.WrapListener(ll)
	}
	return ll, err
}
