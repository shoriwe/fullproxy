package reverse

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/shoriwe/fullproxy/v4/utils/network"
)

type Slave struct {
	initialized bool
	Listener    net.Listener // Optional listener
	Dial        network.DialFunc
	Master      net.Conn       // Master connection
	Control     *yamux.Session // Control channel
}

func (s *Slave) init() (err error) {
	if !s.initialized {
		s.Control, err = yamux.Server(s.Master, yamux.DefaultConfig())
		s.initialized = err == nil
	}
	return err
}

func (s *Slave) HandleAccept(conn net.Conn, req *Request) (err error) {
	if s.Listener == nil {
		err = fmt.Errorf("no listener provided")
	} else {
		var target net.Conn
		target, err = s.Listener.Accept()
		if err == nil {
			defer target.Close()
			err = gob.NewEncoder(conn).Encode(SucceedResponse)
			if err == nil {
				go io.Copy(conn, target)
				_, err = io.Copy(target, conn)
			}
		}
	}
	if err != nil {
		gob.NewEncoder(conn).Encode(FailResponse(err))
	}
	return err
}

func (s *Slave) HandleDial(conn net.Conn, req *Request) (err error) {
	var target net.Conn
	target, err = s.Dial(req.Network, req.Address)
	if err == nil {
		defer target.Close()
		err = gob.NewEncoder(conn).Encode(SucceedResponse)
		if err == nil {
			go io.Copy(conn, target)
			_, err = io.Copy(target, conn)
			return err
		}
	}
	gob.NewEncoder(conn).Encode(FailResponse(err))
	return err
}

func (s *Slave) Handle(conn net.Conn) (err error) {
	defer conn.Close()
	var req Request
	err = gob.NewDecoder(conn).Decode(&req)
	if err == nil {
		switch req.Action {
		case Accept:
			err = s.HandleAccept(conn, &req)
		case Dial:
			err = s.HandleDial(conn, &req)
		default:
			err = fmt.Errorf("invalid action")
			gob.NewEncoder(conn).Encode(FailResponse(err))
		}
	}
	return err
}

func (s *Slave) Serve() (err error) {
	err = s.init()
	if err == nil {
		var conn net.Conn
		for {
			conn, err = s.Control.Accept()
			if err == nil {
				go s.Handle(conn)
			}
		}
	}
	return err
}

func (s *Slave) Close() {
	s.Master.Close()
	if s.Listener != nil {
		s.Listener.Close()
	}
}
