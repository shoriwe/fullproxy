package compose

import "log"

type Compose struct {
	Circuits map[string]*Circuit `yaml:"circuits,omitempty" json:"circuits,omitempty"`
	Proxies  map[string]*Proxy   `yaml:"proxies,omitempty" json:"proxies,omitempty"`
	Slaves   map[string]*Slave   `yaml:"slaves,omitempty" json:"slaves,omitempty"`
}

type service interface {
	Serve() error
}

func startService(errCh chan error, s service) {
	err := s.Serve()
	errCh <- err
}

func startServices[T service](services map[string]T, errCh chan error) {
	for name, target := range services {
		log.Printf("[*] Starting %s\n", name)
		go startService(errCh, target)
	}
}

func (c *Compose) Start() error {
	errCh := make(chan error, len(c.Circuits)+len(c.Proxies)+len(c.Slaves))
	defer close(errCh)
	go startServices(c.Circuits, errCh)
	go startServices(c.Proxies, errCh)
	go startServices(c.Slaves, errCh)
	return <-errCh
}
