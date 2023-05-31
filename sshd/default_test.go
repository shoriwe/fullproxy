package sshd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultClientConfig(t *testing.T) {
	assert.NotNil(t, DefaultClientConfig())
}

func TestDefaultClient(t *testing.T) {
	client := DefaultClient(t)
	defer client.Close()
	assert.NotNil(t, client)
}

func TestKeepAlive(t *testing.T) {
	t.Run("Trigger KeepAlive", func(tt *testing.T) {
		client := DefaultClient(tt)
		defer client.Close()
		go KeepAlive(client)
		time.Sleep(2 * time.Second)
	})
	t.Run("Error KeepAlive", func(tt *testing.T) {
		client := DefaultClient(tt)
		defer client.Close()
		go KeepAlive(client)
		time.Sleep(2 * time.Second)
		client.Close()
		time.Sleep(2 * time.Second)
	})
}
