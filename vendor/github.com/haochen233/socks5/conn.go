package socks5

import (
	"io"
	"net"
	"sync"
	"time"
)

// UDPConn be associated with TCP connections.
// The UDP connection will close immediately, When TCP connection closed,
// UDPConn only use in UDP_ASSOCIATE command.
type UDPConn struct {
	mu        sync.Mutex
	udp       *net.UDPConn
	tcp       *net.TCPConn
	closeChan chan struct{}
}

// NewUDPConn get a *UDPConn through provide a tcp and udp connection.
// the tcp connection is used for socks UDP_ASSOCIATE handshake.
// the udp connection is used for socks udp forwarding.
//
// After UDP_ASSOCIATE handshake, the tcp transit nothing. Its only
// function is udp relay connection still running.
//
// If one of them shuts off, it will close them all.
func NewUDPConn(udp *net.UDPConn, tcp *net.TCPConn) *UDPConn {
	if udp == nil || tcp == nil {
		return nil
	}

	u := &UDPConn{
		udp:       udp,
		tcp:       tcp,
		closeChan: make(chan struct{}),
	}

	go func() {
		// guard tcp connection, if it closed should close tcp relay too.
		io.Copy(io.Discard, tcp)
		u.Close()
	}()

	return u
}

func (u *UDPConn) LocalAddr() net.Addr {
	return u.udp.LocalAddr()
}

func (u *UDPConn) RemoteAddr() net.Addr {
	return u.udp.RemoteAddr()
}

func (u *UDPConn) SetDeadline(t time.Time) error {
	return u.udp.SetDeadline(t)
}

func (u *UDPConn) SetReadDeadline(t time.Time) error {
	return u.udp.SetReadDeadline(t)
}

func (u *UDPConn) SetWriteDeadline(t time.Time) error {
	return u.udp.SetWriteDeadline(t)
}

func (u *UDPConn) Read(b []byte) (n int, err error) {
	return u.udp.Read(b)
}

func (u *UDPConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	return u.udp.WriteToUDP(b, addr)
}

func (u *UDPConn) ReadFromUDP(b []byte) (int, *net.UDPAddr, error) {
	return u.udp.ReadFromUDP(b)
}

func (u *UDPConn) Write(b []byte) (n int, err error) {
	return u.udp.Write(b)
}

func (u *UDPConn) Close() error {
	var err error
	u.mu.Lock()
	defer u.mu.Unlock()

	ch := u.getCloseChanLocked()
	select {
	case <-ch:
		return nil
	default:
		if err2 := u.udp.Close(); err2 != nil {
			err = err2
		}
		if err2 := u.tcp.Close(); err2 != nil {
			err = err2
		}
		close(u.closeChan)
	}
	return err
}

func (u *UDPConn) CloseCh() <-chan struct{} {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.getCloseChanLocked()
}

func (u *UDPConn) getCloseChanLocked() <-chan struct{} {
	if u.closeChan == nil {
		u.closeChan = make(chan struct{})
	}
	return u.closeChan
}
