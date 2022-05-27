package common

import (
	"fmt"
	"io"
	"net"
)

type Sniffer struct {
	ReadSniffer  io.Writer
	WriteSniffer io.Writer
	Connection   net.Conn
}

func (s *Sniffer) Read(p []byte) (n int, err error) {
	if s.ReadSniffer == nil {
		return s.Connection.Read(p)
	}
	length, readReadError := s.Connection.Read(p)
	_, _ = s.ReadSniffer.Write(p[:length])
	_, _ = fmt.Fprintf(s.ReadSniffer, SniffSeparator)
	return length, readReadError
}

func (s *Sniffer) Write(b []byte) (int, error) {
	if s.WriteSniffer == nil {
		return s.Connection.Write(b)
	}
	length, readReadError := s.Connection.Write(b)
	_, _ = s.WriteSniffer.Write(b[:length])
	_, _ = fmt.Fprintf(s.WriteSniffer, SniffSeparator)
	return length, readReadError
}

func (s *Sniffer) Close() error {
	return s.Connection.Close()
}
