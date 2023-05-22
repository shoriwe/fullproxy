package config

import (
	"net"

	"github.com/shoriwe/fullproxy/v3/reverse"
)

type Reverse struct {
	Listener   *Listener `yaml:"listener"`
	Controller Listener  `yaml:"controller"`
}

func (r *Reverse) Master() (net.Listener, error) {
	l, lErr := r.Listener.Listen()
	if lErr != nil {
		return nil, lErr
	}
	cl, clErr := r.Listener.Listen()
	if clErr != nil {
		l.Close()
		return nil, clErr
	}
	return reverse.NewMaster(l, cl)
}

func (r *Reverse) Slave() (*reverse.Slave, error) {
	cl, clErr := r.Controller.Dial()
	if clErr != nil {
		return nil, clErr
	}
	return reverse.NewSlave(cl)
}
