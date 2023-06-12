package proxies

import (
	"net"
	"sync"
	"testing"

	"github.com/shoriwe/fullproxy/v4/reverse"
	"github.com/shoriwe/fullproxy/v4/utils/network"
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
	// - Proxy
	l := network.ListenAny()
	defer l.Close()
	// - Service
	service := network.ListenAny()
	defer service.Close()
	// Proxy
	s := Socks5{
		Listener: l,
		Dial:     net.Dial,
	}
	go s.Serve()
	// Test
	testMsg := []byte("HELLO")
	// - Producer
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := service.Accept()
		assert.Nil(t, err)
		defer conn.Close()
		// Write
		_, err = conn.Write(testMsg)
		assert.Nil(t, err)
		// Read
		buffer := make([]byte, len(testMsg))
		_, err = conn.Read(buffer)
		assert.Nil(t, err)
	}()
	// - Consumer
	// -- Connect proxy
	dialer, err := proxy.SOCKS5(l.Addr().Network(), l.Addr().String(), nil, nil)
	assert.Nil(t, err)
	// -- Dial service
	conn, err := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, err)
	defer conn.Close()
	// Consume
	buffer := make([]byte, len(testMsg))
	_, err = conn.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, testMsg, buffer)
	// -- Send
	_, err = conn.Write(buffer)
	assert.Nil(t, err)
}

func TestSocks5_Reverse(t *testing.T) {
	// - Proxy
	data := network.ListenAny()
	defer data.Close()
	control := network.ListenAny()
	defer control.Close()
	masterConn := network.Dial(control.Addr().String())
	defer masterConn.Close()
	// - Service
	service := network.ListenAny()
	defer service.Close()
	// - Slave
	slave := &reverse.Slave{
		Master: masterConn,
		Dial:   net.Dial,
	}
	defer slave.Close()
	go slave.Serve()
	// - Master
	master := &reverse.Master{
		Data:    data,
		Control: control,
	}
	defer master.Close()
	// - Setup Proxy
	sockProxy := Socks5{
		Listener: master,
		Dial:     master.SlaveDial,
	}
	go sockProxy.Serve()
	// - Producer
	var wg sync.WaitGroup
	defer wg.Wait()
	testMsg := []byte("HELLO")
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := service.Accept()
		assert.Nil(t, err)
		defer conn.Close()
		// Write
		_, err = conn.Write(testMsg)
		assert.Nil(t, err)
		buffer := make([]byte, len(testMsg))
		// Read
		_, err = conn.Read(buffer)
		assert.Nil(t, err)
	}()
	// - Consumer
	// -- Connect proxy
	dialer, err := proxy.SOCKS5(data.Addr().Network(), data.Addr().String(), nil, nil)
	assert.Nil(t, err)
	// -- Connect service
	conn, err := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, err)
	defer conn.Close()
	// -- Consume
	buffer := make([]byte, len(testMsg))
	_, err = conn.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, testMsg, buffer)
	// -- Produce
	_, err = conn.Write(buffer)
	assert.Nil(t, err)
}

func TestSocks5_UsernamePassword(t *testing.T) {
	// Proxy
	l := network.ListenAny()
	defer l.Close()
	// Service
	service := network.ListenAny()
	defer service.Close()
	// Setup Proxy
	s := Socks5{
		Listener: l,
		Dial:     net.Dial,
		AuthMethods: []socks5.Authenticator{
			socks5.UserPassAuthenticator{
				Credentials: socks5.StaticCredentials{"username": "password"},
			},
		},
	}
	defer s.Close()
	go s.Serve()
	// Producer
	testMsg := []byte("HELLO")
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := service.Accept()
		assert.Nil(t, err)
		defer conn.Close()
		// Write
		_, err = conn.Write(testMsg)
		assert.Nil(t, err)
		// Read
		buffer := make([]byte, len(testMsg))
		_, err = conn.Read(buffer)
		assert.Nil(t, err)
	}()
	// Consumer
	// Connect proxy
	dialer, err := proxy.SOCKS5(l.Addr().Network(), l.Addr().String(), &proxy.Auth{
		User:     "username",
		Password: "password",
	}, nil)
	assert.Nil(t, err)
	// Connect service
	conn, dErr := dialer.Dial(service.Addr().Network(), service.Addr().String())
	assert.Nil(t, dErr)
	defer conn.Close()
	// Consume
	buffer := make([]byte, len(testMsg))
	_, err = conn.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, testMsg, buffer)
	// Produce
	_, err = conn.Write(buffer)
	assert.Nil(t, err)
}
