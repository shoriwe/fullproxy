package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := ListenAny()
		defer l.Close()
		go func() {
			conn, err := l.Accept()
			assert.Nil(tt, err)
			defer conn.Close()
		}()
		conn := Dial(l.Addr().String())
		defer conn.Close()
	})
}
