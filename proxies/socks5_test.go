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
	data := network.ListenAny()
	defer data.Close()
	control := network.ListenAny()
	defer control.Close()
	master := network.Dial(control.Addr().String())
	defer master.Close()
	doneChan := make(chan struct{}, 1)
	defer close(doneChan)
	go func() {
		s := &reverse.Slave{
			Master: master,
		}
		defer s.Close()
		go s.Serve()
		<-doneChan
	}()
	m := &reverse.Master{
		Data:    data,
		Control: control,
	}
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
		Listener: m,
		Dial:     m.SlaveDial,
	}
	go s.Serve()
	dialer, sErr := proxy.SOCKS5(data.Addr().Network(), data.Addr().String(), nil, nil)
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
