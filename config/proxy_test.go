package config

import (
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
)

func TestProxy_Create(t *testing.T) {
	t.Run("Basic Listener", func(tt *testing.T) {
		p := Proxy{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
			Protocol: Socks5,
		}
		proxy, pErr := p.Create()
		assert.Nil(tt, pErr)
		defer proxy.Close()
	})
	t.Run("Reverse Listener", func(tt *testing.T) {
		p := Proxy{
			Reverse: &Reverse{
				Listener: &Listener{
					Network: "tcp",
					Address: "localhost:9999",
				},
				Controller: Listener{
					Network: "tcp",
					Address: "localhost:10000",
				},
			},
			Protocol: Socks5,
		}
		doneCh := make(chan struct{}, 1)
		go func() {
			r := Reverse{
				Controller: Listener{
					Network: "tcp",
					Address: "localhost:10000",
				},
			}
			slave, sErr := r.Slave()
			assert.Nil(tt, sErr)
			defer slave.Close()
			go slave.Serve()
			<-doneCh
		}()
		proxy, pErr := p.Create()
		assert.Nil(tt, pErr)
		defer proxy.Close()
	})
	t.Run("Invalid Protocol", func(tt *testing.T) {
		p := Proxy{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
			Protocol: "INVALID",
		}
		_, pErr := p.Create()
		assert.NotNil(tt, pErr)
	})
	t.Run("No Listeners", func(tt *testing.T) {
		p := Proxy{}
		_, pErr := p.Create()
		assert.NotNil(tt, pErr)
	})
	t.Run("Socks5", func(tt *testing.T) {
		p := Proxy{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
			Protocol: Socks5,
		}
		px, pErr := p.Create()
		assert.Nil(tt, pErr)
		defer px.Close()
		go px.Serve()
		// Test the proxy
		service := network.ListenAny()
		defer service.Close()
		msg := []byte("TEST")
		doneCh := make(chan struct{}, 1)
		go func() {
			conn, aErr := service.Accept()
			assert.Nil(tt, aErr)
			defer conn.Close()
			_, wErr := conn.Write(msg)
			assert.Nil(tt, wErr)
			<-doneCh
		}()
		dialer, dErr := proxy.SOCKS5(px.Addr().Network(), px.Addr().String(), nil, nil)
		assert.Nil(tt, dErr)
		conn, cErr := dialer.Dial(service.Addr().Network(), service.Addr().String())
		assert.Nil(tt, cErr)
		buffer := make([]byte, len(msg))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		doneCh <- struct{}{}
		assert.Equal(tt, msg, buffer)
	})
}
