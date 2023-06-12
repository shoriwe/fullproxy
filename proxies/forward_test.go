package proxies

import (
	"net"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestForward_Addr(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	service := network.ListenAny()
	defer service.Close()
	f := Forward{
		Network:  service.Addr().Network(),
		Address:  service.Addr().String(),
		Listener: listener,
		Dial:     net.Dial,
	}
	defer f.Close()
	go f.Serve()
	assert.NotNil(t, f.Addr())
	testMessage := []byte("TEST")
	go func() {
		conn, err := service.Accept()
		assert.Nil(t, err)
		defer conn.Close()
		// Write
		_, err = conn.Write(testMessage)
		assert.Nil(t, err)
		// Read
		buffer := make([]byte, len(testMessage))
		_, err = conn.Read(buffer)
		assert.Nil(t, err)
		assert.Equal(t, testMessage, buffer)
	}()
	conn := network.Dial(listener.Addr().String())
	defer conn.Close()
	// Read
	buffer := make([]byte, len(testMessage))
	_, err := conn.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, testMessage, buffer)
	// Write
	_, err = conn.Write(buffer)
	assert.Nil(t, err)
}

func TestBasicLocalForward(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	service := network.ListenAny()
	defer service.Close()
	f := Forward{
		Network:  service.Addr().Network(),
		Address:  service.Addr().String(),
		Listener: listener,
		Dial:     net.Dial,
	}
	defer f.Close()
	go f.Serve()
	testMessage := []byte("TEST")
	go func() {
		conn, err := service.Accept()
		assert.Nil(t, err)
		defer conn.Close()
		// Write
		_, err = conn.Write(testMessage)
		assert.Nil(t, err)
		// Read
		buffer := make([]byte, len(testMessage))
		_, err = conn.Read(buffer)
		assert.Nil(t, err)
		assert.Equal(t, testMessage, buffer)
	}()
	conn := network.Dial(listener.Addr().String())
	defer conn.Close()
	// Read
	buffer := make([]byte, len(testMessage))
	_, err := conn.Read(buffer)
	assert.Nil(t, err)
	assert.Equal(t, testMessage, buffer)
	// Write
	_, err = conn.Write(testMessage)
	assert.Nil(t, err)

}
