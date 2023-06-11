package compose

import (
	"fmt"
	"io"
	"net"

	"github.com/shoriwe/fullproxy/v3/circuit"
)

type Circuit struct {
	Network  string  `yaml:"network" json:"network"`
	Address  string  `yaml:"address" json:"address"`
	Listener Network `yaml:"listener" json:"listener"`
	Knots    []Knot  `yaml:"knots" json:"knots"`
	listener net.Listener
	circuit  *circuit.Circuit
}

func (c *Circuit) handle(conn net.Conn) {
	defer conn.Close()
	target, dErr := c.circuit.Dial(c.Network, c.Address)
	if dErr == nil {
		defer target.Close()
		go io.Copy(conn, target)
		io.Copy(target, conn)
	}
}

func (c *Circuit) setupCircuit() error {
	c.circuit = &circuit.Circuit{Chain: make([]circuit.Knot, 0, len(c.Knots))}
	for index, k := range c.Knots {
		knot, err := k.Compile()
		if err != nil {
			return fmt.Errorf("error compiling knot %d", index+1)
		}
		c.circuit.Chain = append(c.circuit.Chain, knot)
	}
	return nil
}

func (c *Circuit) serve(l net.Listener) (err error) {
	var conn net.Conn
	for conn, err = l.Accept(); err == nil; conn, err = l.Accept() {
		go c.handle(conn)
	}
	return err
}

func (c *Circuit) Serve() (err error) {
	err = c.setupCircuit()
	if err == nil {
		c.listener, err = c.Listener.Listen()
		if err == nil {
			defer c.listener.Close()
			return c.serve(c.listener)
		}
	}
	return err
}
