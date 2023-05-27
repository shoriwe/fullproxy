package reverse

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/hashicorp/yamux"
)

type Master struct {
	net.Listener
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

func (m *Master) Dial(network, addr string) (net.Conn, error) {
	stream, sErr := m.cSession.Open()
	if sErr != nil {
		return nil, sErr
	}
	eErr := gob.NewEncoder(stream).Encode(Request{
		Action:  Dial,
		Network: network,
		Address: addr,
	})
	if eErr != nil {
		stream.Close()
		return nil, eErr
	}
	var response Response
	dErr := gob.NewDecoder(stream).Decode(&response)
	if dErr != nil {
		stream.Close()
		return nil, dErr
	}
	if response.Succeed {
		return stream, nil
	}
	stream.Close()
	return nil, fmt.Errorf(response.Message)
}

func (m *Master) Accept() (net.Conn, error) {
	return m.Listener.Accept()
}

func (m *Master) Close() error {
	m.Listener.Close()
	m.cConn.Close()
	m.cListener.Close()
	return nil
}

func NewMaster(listener, controlListener net.Listener) (*Master, error) {
	m := &Master{
		Listener:  listener,
		cListener: controlListener,
	}
	iErr := m.init()
	return m, iErr
}
