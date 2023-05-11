package listeners

import (
	"crypto/tls"
	"net"

	"github.com/hashicorp/yamux"
)

type Master struct {
	listener  net.Listener
	cListener net.Listener   // Control listener
	cConn     net.Conn       // Control connection
	cSession  *yamux.Session // Control session
}

func (m *Master) init() error {
	var err error
	m.cConn, err = m.cListener.Accept()
	if err != nil {
		return err
	}
	m.cSession, err = yamux.Client(m.cConn, yamux.DefaultConfig())
	return err
}

func (m *Master) Accept() (net.Conn, error) {
	session, sErr := m.cSession.Open()
	if sErr != nil {
		return nil, sErr
	}
	ses
}

func NewMaster(addr, cAddr string, cConfig *tls.Config) (net.Listener, error) {
	l, lErr := net.Listen("tcp", addr)
	if lErr != nil {
		return nil, lErr
	}
	cl, clErr := tls.Listen("tcp", cAddr, cConfig)
	if clErr != nil {
		return nil, clErr
	}
	m := &Master{
		listener:  l,
		cListener: cl,
	}
	m.init()
	return m, nil
}
