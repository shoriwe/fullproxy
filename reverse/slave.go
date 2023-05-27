package reverse

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/hashicorp/yamux"
)

type Slave struct {
	initialized bool
	Listener    net.Listener   // Optional listener
	Control     net.Conn       // Control channel
	Data        *yamux.Session // Data channel
}

func (s *Slave) init() error {
	if s.initialized {
		return nil
	}
	var err error
	s.Data, err = yamux.Server(s.Control, yamux.DefaultConfig())
	s.initialized = true
	return err
}

func (s *Slave) handleAccept(stream net.Conn, req *Request) {
	target, aErr := s.Listener.Accept()
	if aErr == nil {
		defer target.Close()
		gob.NewEncoder(stream).Encode(SucceedResponse)
		go io.Copy(stream, target)
		io.Copy(target, stream)
		return
	}
	gob.NewEncoder(stream).Encode(FailResponse(aErr))
}

func (s *Slave) handleDial(stream net.Conn, req *Request) {
	target, dialErr := net.Dial(req.Network, req.Address)
	if dialErr == nil {
		defer target.Close()
		gob.NewEncoder(stream).Encode(SucceedResponse)
		go io.Copy(stream, target)
		io.Copy(target, stream)
		return
	}
	gob.NewEncoder(stream).Encode(FailResponse(dialErr))
}

func (s *Slave) Handle(stream net.Conn) {
	defer stream.Close()
	var req Request
	gob.NewDecoder(stream).Decode(&req)
	switch req.Action {
	case Accept:
		s.handleAccept(stream, &req)
		return
	case Dial:
		s.handleDial(stream, &req)
		return
	default:
		gob.NewEncoder(stream).Encode(FailResponse(fmt.Errorf("invalid action")))
		return
	}
}

func (s *Slave) Serve() error {
	iErr := s.init()
	if iErr != nil {
		return iErr
	}
	for {
		stream, err := s.Data.Accept()
		if err == nil {
			go s.Handle(stream)
		}
	}
}
func (s *Slave) Close() {
	s.Control.Close()
}
