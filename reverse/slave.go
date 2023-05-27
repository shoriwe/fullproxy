package reverse

import (
	"encoding/gob"
	"io"
	"net"

	"github.com/hashicorp/yamux"
)

type Slave struct {
	cConn    net.Conn       // Control connection
	cSession *yamux.Session // Control session
}

func (s *Slave) init() error {
	var err error
	s.cSession, err = yamux.Server(s.cConn, yamux.DefaultConfig())
	return err
}

func (s *Slave) Handle(stream net.Conn) {
	defer stream.Close()
	var addr string
	dErr := gob.NewDecoder(stream).Decode(&addr)
	if dErr != nil {
		return // TODO: Log error?
	}
	target, dialErr := net.Dial("tcp", addr)
	if dialErr != nil {
		gob.NewEncoder(stream).Encode(Response{Succeed: false, Message: dialErr.Error()})
		return // TODO: Log error?
	}
	defer target.Close()
	eErr := gob.NewEncoder(stream).Encode(Response{Succeed: true, Message: "Succeed"})
	if eErr != nil {
		return // TODO: Log error?
	}
	go io.Copy(stream, target)
	io.Copy(target, stream) // TODO: Log error?
}

func (s *Slave) Serve() {
	for {
		stream, err := s.cSession.Accept()
		if err == nil {
			go s.Handle(stream)
		}
	}
}
func (s *Slave) Close() {
	s.cConn.Close()
}

func NewSlave(conn net.Conn) (*Slave, error) {
	s := &Slave{
		cConn: conn,
	}
	iErr := s.init()
	return s, iErr
}
