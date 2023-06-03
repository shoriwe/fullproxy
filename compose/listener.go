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
	ListenerBasic  = "basic"
	ListenerMaster = "master"
	ListenerSSH    = "ssh"
)

type Listener struct {
	Type    string    `yaml:"type" json:"type"`
	Network *string   `yaml:"network,omitempty" json:"network,omitempty"`
	Address *string   `yaml:"address,omitempty" json:"address,omitempty"`
	Data    *Listener `yaml:"data,omitempty" json:"data,omitempty"`
	Control *Listener `yaml:"control,omitempty" json:"control,omitempty"`
	Auth    *Auth     `yaml:"auth,omitempty" json:"auth,omitempty"`
	Crypto  *Crypto   `yaml:"crypto,omitempty" json:"crypto,omitempty"`
}

func (l *Listener) setupBasic(listen network.ListenFunc) (net.Listener, error) {
	if l.Network == nil {
		return nil, fmt.Errorf("network not set for basic listener")
	}
	if l.Address == nil {
		return nil, fmt.Errorf("address not set for basic listener")
	}
	return listen(*l.Network, *l.Address)
}

func (l *Listener) setupMaster() (ll net.Listener, err error) {
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

func (l *Listener) setupSSH() (ll net.Listener, err error) {
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
	data, err = l.Data.setupBasic(sshConn.Listen)
	if err == nil {
		ll = &sshWrapper{
			Listener: data,
			conn:     sshConn,
		}
	}
	return ll, err
}

func (l *Listener) Listen() (ll net.Listener, err error) {
	switch l.Type {
	case ListenerBasic:
		ll, err = l.setupBasic(net.Listen)
	case ListenerMaster:
		ll, err = l.setupMaster()
	case ListenerSSH:
		ll, err = l.setupSSH()
	default:
		err = fmt.Errorf("unknown listener type %s", l.Type)
	}
	if err == nil && l.Crypto != nil {
		ll, err = l.Crypto.WrapListener(ll)
	}
	return ll, err
}
