package compose

import (
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
)

func Test_startServices(t *testing.T) {
	service := network.ListenAny()
	defer service.Close()
	exampleName := "example"
	p := &Proxy{
		Type:    ProxyForward,
		Network: new(string),
		Address: new(string),
		Listener: Network{
			Type:    NetworkBasic,
			Network: new(string),
			Address: new(string),
		},
	}
	*p.Network = service.Addr().Network()
	*p.Address = service.Addr().String()
	*p.Listener.Network = "tcp"
	*p.Listener.Address = "localhost:0"
	c := Compose{
		Proxies: map[string]*Proxy{
			exampleName: p,
		},
	}
	errCh := make(chan error, 1)
	go startServices(c.Proxies, errCh)
	for p.Listener.listener == nil {
	}
	defer p.Listener.listener.Close()
}

func TestCompose_Start(t *testing.T) {
	service := network.ListenAny()
	defer service.Close()
	exampleName := "example"
	p := &Proxy{
		Type:    ProxyForward,
		Network: new(string),
		Address: new(string),
		Listener: Network{
			Type:    NetworkBasic,
			Network: new(string),
			Address: new(string),
		},
	}
	*p.Network = service.Addr().Network()
	*p.Address = service.Addr().String()
	*p.Listener.Network = "tcp"
	*p.Listener.Address = "localhost:0"
	c := Compose{
		Proxies: map[string]*Proxy{
			exampleName: p,
		},
	}
	go c.Start()
	for p.Listener.listener == nil {
	}
	defer p.Listener.listener.Close()
}
