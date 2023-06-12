package compose

import (
	"fmt"

	"github.com/shoriwe/fullproxy/v4/circuit"
	"golang.org/x/net/proxy"
)

const (
	KnotForward = "forward"
	KnotSocks5  = "socks5"
	KnotSSH     = "ssh"
)

var AvailableKnots = [3]string{KnotForward, KnotSocks5, KnotSSH}

type Knot struct {
	Type    string `yaml:"type,omitempty" json:"type,omitempty"`
	Network string `yaml:"network,omitempty" json:"network,omitempty"`
	Address string `yaml:"address,omitempty" json:"address,omitempty"`
	Auth    *Auth  `yaml:"auth,omitempty" json:"auth,omitempty"`
}

func (k *Knot) Compile() (circuit.Knot, error) {
	switch k.Type {
	case KnotForward:
		return &circuit.Forward{
			Network: k.Network,
			Address: k.Address,
		}, nil
	case KnotSocks5:
		var (
			auth *proxy.Auth
			err  error
		)
		if k.Auth != nil {
			auth, err = k.Auth.Socks5()
			if err != nil {
				return nil, err
			}
		}
		return &circuit.Socks5{
			Network: k.Network,
			Address: k.Address,
			Auth:    auth,
		}, nil
	case KnotSSH:
		if k.Auth == nil {
			return nil, fmt.Errorf("no auth provided for SSH")
		}
		config, err := k.Auth.SSHClientConfig()
		if err != nil {
			return nil, err
		}
		return &circuit.SSH{
			Network: k.Network,
			Address: k.Address,
			Config:  *config,
		}, nil
	default:
		return nil, fmt.Errorf("unknown knot %s; available knots are %s", k.Type, AvailableKnots)
	}
}
