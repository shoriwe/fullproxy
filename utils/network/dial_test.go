package network

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
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
	t.Run("Invalid", func(tt *testing.T) {
		defer func() {
			assert.NotNil(tt, recover())
		}()
		conn := Dial("11111111111111111")
		defer conn.Close()
	})
}

func TestNopClose(t *testing.T) {
	assert.Nil(t, NopClose())
}

func TestCloseOnError(t *testing.T) {
	listener := ListenAny()
	defer listener.Close()
	go listener.Accept()
	conn := Dial(listener.Addr().String())
	defer conn.Close()
	err := fmt.Errorf("an error")
	CloseOnError(&err, conn)
}
