package filter

import (
	"fmt"
	"net"
)

type Listener struct {
	net.Listener
	Whitelist []Match
	Blacklist []Match
}

func (l *Listener) GetWhitelist() []Match { return l.Whitelist }
func (l *Listener) GetBlacklist() []Match { return l.Blacklist }

func (l *Listener) Accept() (conn net.Conn, err error) {
	conn, err = l.Listener.Accept()
	if err == nil {
		if !VerifyConn(l, conn.RemoteAddr()) {
			conn.Close()
			err = fmt.Errorf("connection denied by rule")
		}
	}
	return conn, err
}
