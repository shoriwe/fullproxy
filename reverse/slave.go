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

func (s *Slave) Handle() error {
	stream, err := s.cSession.Accept()
	if err != nil {
		return err
	}
	var addr string
	dErr := gob.NewDecoder(stream).Decode(&addr)
	if dErr != nil {
		stream.Close()
		return dErr
	}
	target, dialErr := net.Dial("tcp", addr)
	if dialErr != nil {
		gob.NewEncoder(stream).Encode(Response{Succeed: false, Message: dialErr.Error()})
		stream.Close()
		return dialErr
	}
	eErr := gob.NewEncoder(stream).Encode(Response{Succeed: true, Message: "Succeed"})
	if eErr == nil {
		go io.Copy(stream, target)
		_, cErr := io.Copy(target, stream)
		return cErr
	}
	stream.Close()
	return eErr
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
