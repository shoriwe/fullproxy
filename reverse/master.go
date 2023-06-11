package reverse

import (
	"encoding/gob"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

type Master struct {
	Data        net.Listener
	Control     net.Listener   // Control listener
	Slave       net.Conn       // Slave connection
	cSession    *yamux.Session // Control session
	initialized bool
}

func (m *Master) init() (err error) {
	if !m.initialized {
		m.Slave, err = m.Control.Accept()
		if err == nil {
			m.cSession, err = yamux.Client(m.Slave, yamux.DefaultConfig())
			m.initialized = err == nil
		}
	}
	return err
}

func (m *Master) handle(req *Request) (conn net.Conn, err error) {
	err = m.init()
	if err == nil {
		conn, err = m.cSession.Open()
		if err == nil {
			defer network.CloseOnError(&err, conn)
			err = gob.NewEncoder(conn).Encode(req)
			if err == nil {
				var response Response
				err = gob.NewDecoder(conn).Decode(&response)
				if err == nil {
					err = response.Message
				}
			}
		}
	}
	return conn, err
}

func (m *Master) SlaveAccept() (net.Conn, error) {
	req := Request{
		Action: Accept,
	}
	return m.handle(&req)
}

func (m *Master) SlaveDial(network, addr string) (net.Conn, error) {
	req := Request{
		Action:  Dial,
		Network: network,
		Address: addr,
	}
	return m.handle(&req)
}

func (m *Master) Accept() (net.Conn, error) {
	return m.Data.Accept()
}

func (m *Master) Addr() net.Addr {
	return m.Data.Addr()
}

func (m *Master) Close() error {
	m.Data.Close()
	m.Control.Close()
	if m.Slave != nil {
		m.Slave.Close()
	}
	return nil
}
