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
	doneChan := make(chan struct{}, 1)
	defer close(doneChan)
	testMessage := []byte("TEST")
	go func() {
		conn, aErr := service.Accept()
		assert.Nil(t, aErr)
		defer conn.Close()
		_, wErr := conn.Write(testMessage)
		assert.Nil(t, wErr)
		<-doneChan
	}()
	go f.Serve()
	conn, dErr := net.Dial(listener.Addr().Network(), listener.Addr().String())
	assert.Nil(t, dErr)
	defer conn.Close()
	buffer := make([]byte, len(testMessage))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, testMessage, buffer)
	assert.NotNil(t, f.Addr())
	doneChan <- struct{}{}
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
	doneChan := make(chan struct{}, 1)
	defer close(doneChan)
	testMessage := []byte("TEST")
	go func() {
		conn, aErr := service.Accept()
		assert.Nil(t, aErr)
		defer conn.Close()
		_, wErr := conn.Write(testMessage)
		assert.Nil(t, wErr)
		<-doneChan
	}()
	go f.Serve()
	conn, dErr := net.Dial(listener.Addr().Network(), listener.Addr().String())
	assert.Nil(t, dErr)
	defer conn.Close()
	buffer := make([]byte, len(testMessage))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, testMessage, buffer)
	doneChan <- struct{}{}

}
