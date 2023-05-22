package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListener_Listen(t *testing.T) {
	t.Run("No TLS", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:0",
		}
		l, lErr := c.Listen()
		assert.Nil(tt, lErr)
		defer l.Close()
	})
	t.Run("With TLS", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:0",
			TLS:     &TLS{},
		}
		l, lErr := c.Listen()
		assert.Nil(tt, lErr)
		defer l.Close()
	})
	t.Run("Invalid Addr", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:9999999",
			TLS:     &TLS{},
		}
		_, lErr := c.Listen()
		assert.NotNil(tt, lErr)
	})
	t.Run("Invalid TLS", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:0",
			TLS: &TLS{
				CertFile: new(string),
				KeyFile:  new(string),
			},
		}
		_, lErr := c.Listen()
		assert.NotNil(tt, lErr)
	})
}

func TestListener_Dial(t *testing.T) {
	t.Run("No TLS", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:0",
		}
		l, lErr := c.Listen()
		assert.Nil(tt, lErr)
		defer l.Close()
		doneCh := make(chan struct{}, 1)
		msg := []byte("TEST")
		go func() {
			c2 := Listener{
				Network: l.Addr().Network(),
				Address: l.Addr().String(),
			}
			conn, cErr := c2.Dial()
			assert.Nil(tt, cErr)
			defer conn.Close()
			_, wErr := conn.Write(msg)
			assert.Nil(tt, wErr)
			<-doneCh
		}()
		conn, aErr := l.Accept()
		assert.Nil(tt, aErr)
		buffer := make([]byte, len(msg))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, msg, buffer)
		doneCh <- struct{}{}
	})
	t.Run("With TLS", func(tt *testing.T) {
		c := Listener{
			Network: "tcp",
			Address: "localhost:0",
			TLS:     &TLS{},
		}
		l, lErr := c.Listen()
		assert.Nil(tt, lErr)
		defer l.Close()
		doneCh := make(chan struct{}, 1)
		msg := []byte("TEST")
		go func() {
			c2 := Listener{
				Network: l.Addr().Network(),
				Address: l.Addr().String(),
				TLS:     &TLS{},
			}
			conn, cErr := c2.Dial()
			assert.Nil(tt, cErr)
			defer conn.Close()
			_, wErr := conn.Write(msg)
			assert.Nil(tt, wErr)
			<-doneCh
		}()
		conn, aErr := l.Accept()
		assert.Nil(tt, aErr)
		buffer := make([]byte, len(msg))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, msg, buffer)
		doneCh <- struct{}{}
	})

}
