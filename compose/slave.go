package compose

import (
	"net"

	"github.com/shoriwe/fullproxy/v4/reverse"
	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type Slave struct {
	MasterNetwork string   `yaml:"masterNetwork" json:"masterNetwork"`
	MasterAddress string   `yaml:"masterAddress" json:"masterAddress"`
	MasterDialer  Network  `yaml:"masterDialer" json:"masterDialer"`
	Dialer        Network  `yaml:"dialer" json:"dialer"`
	Listener      *Network `yaml:"listener,omitempty" json:"listener,omitempty"`
	slave         *reverse.Slave
}

func (s *Slave) setupSlave() (*reverse.Slave, error) {
	if s.slave == nil {
		masterDialFunc, err := s.MasterDialer.DialFunc()
		if err == nil {
			masterConn, err := masterDialFunc(s.MasterNetwork, s.MasterAddress)
			if err == nil {
				network.CloseOnError(&err, masterConn)
				dialFunc, err := s.Dialer.DialFunc()
				if err == nil {
					var l net.Listener
					if s.Listener != nil {
						l, err = s.Listener.Listen()
						network.CloseOnError(&err, l)
					}
					if err == nil {
						s.slave = &reverse.Slave{
							Listener: l,
							Dial:     dialFunc,
							Master:   masterConn,
						}
					}
				}
			}
		}
	}
	return s.slave, nil
}

func (s *Slave) Serve() (err error) {
	var slave *reverse.Slave
	slave, err = s.setupSlave()
	if err == nil {
		err = slave.Serve()
	}
	return err
}
