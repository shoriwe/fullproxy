package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListener(t *testing.T) {
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
