package main

import (
	"fmt"
	"github.com/shoriwe/fullproxy/v3/internal/listeners"
	"io"
	"net"
)

type hijackConn struct {
	net.Conn
	read  func(b []byte) (int, error)
	write func(b []byte) (int, error)
}

func (h *hijackConn) Read(b []byte) (n int, err error) {
	return h.read(b)
}

func (h *hijackConn) Write(b []byte) (n int, err error) {
	return h.write(b)
}

type hijackListener struct {
	listeners.Listener
	incoming, outgoing io.Writer
}

func (h *hijackListener) Dial(networkType, address string) (net.Conn, error) {
	conn, connectionError := h.Listener.Dial(networkType, address)
	if connectionError != nil {
		return nil, connectionError
	}
	result := &hijackConn{
		Conn: conn,
	}
	if h.incoming != nil {
		result.read = func(b []byte) (int, error) {
			length, readError := conn.Read(b)
			_, _ = fmt.Fprintf(h.incoming, "\n\n--------------------------------\n\n")
			_, _ = h.incoming.Write(b[:length])
			return length, readError
		}
	} else {
		result.read = conn.Read
	}
	if h.outgoing != nil {
		result.write = func(b []byte) (int, error) {
			_, _ = h.outgoing.Write(b)
			return conn.Write(b)
		}
	} else {
		result.read = conn.Write
	}
	return result, nil
}
