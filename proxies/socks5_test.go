package proxies

import (
	"net"
	"testing"

	"github.com/shoriwe/fullproxy/v3/reverse"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
	"github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

func TestSocks5_Addr(t *testing.T) {
	l := network.ListenAny()
	defer l.Close()
	s := Socks5{
		Listener: l,
		Dial:     net.Dial,
	}
	defer s.Close()
	assert.NotNil(t, s.Addr())
}

func TestSocks5_Listener(t *testing.T) {
	l := network.ListenAny()
	defer l.Close()
	service := network.ListenAny()
	defer service.Close()
	msg := []byte("HELLO")
	go func() {
		conn, aErr := service.Accept()
		assert.Nil(t, aErr)
		_, wErr := conn.Write(msg)
		assert.Nil(t, wErr)
	}()
	s := Socks5{
		Listener: l,
		Dial:     net.Dial,
	}
	go s.Serve()
	dialer, sErr := proxy.SOCKS5(l.Addr().Network(), l.Addr().String(), nil, nil)
	assert.Nil(t, sErr)
	conn, dErr := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, dErr)
	buffer := make([]byte, len(msg))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, msg, buffer)
}

func TestSocks5_Reverse(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	controlListener := network.ListenAny()
	defer controlListener.Close()
	slaveConn := network.Dial(controlListener.Addr().String())
	defer slaveConn.Close()
	doneChan := make(chan struct{}, 1)
	defer close(doneChan)
	go func() {
		slave, err := reverse.NewSlave(slaveConn)
		assert.Nil(t, err)
		defer slave.Close()
		go slave.Handle()
		<-doneChan
	}()
	master, mErr := reverse.NewMaster(listener, controlListener)
	assert.Nil(t, mErr)
	service := network.ListenAny()
	defer service.Close()
	msg := []byte("HELLO")
	go func() {
		conn, aErr := service.Accept()
		assert.Nil(t, aErr)
		_, wErr := conn.Write(msg)
		assert.Nil(t, wErr)
	}()
	s := Socks5{
		Listener: master,
		Dial:     master.Dial,
	}
	go s.Serve()
	dialer, sErr := proxy.SOCKS5(listener.Addr().Network(), listener.Addr().String(), nil, nil)
	assert.Nil(t, sErr)
	conn, dErr := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, dErr)
	buffer := make([]byte, len(msg))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, msg, buffer)
	doneChan <- struct{}{}
}

func TestSocks5_UsernamePassword(t *testing.T) {
	l := network.ListenAny()
	defer l.Close()
	service := network.ListenAny()
	defer service.Close()
	msg := []byte("HELLO")
	go func() {
		conn, aErr := service.Accept()
		assert.Nil(t, aErr)
		_, wErr := conn.Write(msg)
		assert.Nil(t, wErr)
	}()
	s := Socks5{
		Listener: l,
		Dial:     net.Dial,
		AuthMethods: []socks5.Authenticator{
			socks5.UserPassAuthenticator{
				Credentials: socks5.StaticCredentials{"username": "password"},
			},
		},
	}
	go s.Serve()
	dialer, sErr := proxy.SOCKS5(l.Addr().Network(), l.Addr().String(), &proxy.Auth{
		User:     "username",
		Password: "password",
	}, nil)
	assert.Nil(t, sErr)
	conn, dErr := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, dErr)
	buffer := make([]byte, len(msg))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, msg, buffer)
}
