package reverse

import (
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

const (
	testMessage = "MESSAGE"
)

func TestNewMaster(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		listener := network.ListenAny()
		defer listener.Close()
		controlListener := network.ListenAny()
		defer controlListener.Close()
		slaveConn := network.Dial(controlListener.Addr().String())
		defer slaveConn.Close()
		doneChan := make(chan struct{}, 1)
		defer close(doneChan)
		go func() {
			slave, err := NewSlave(slaveConn)
			assert.Nil(tt, err)
			defer slave.Close()
			<-doneChan
		}()
		master, mErr := NewMaster(listener, controlListener)
		assert.Nil(tt, mErr)
		defer master.Close()
		doneChan <- struct{}{}
	})
}

func TestMaster_Accept(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		listener := network.ListenAny()
		defer listener.Close()
		controlListener := network.ListenAny()
		defer controlListener.Close()
		slaveConn := network.Dial(controlListener.Addr().String())
		defer slaveConn.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		go func() {
			slave, err := NewSlave(slaveConn)
			assert.Nil(tt, err)
			defer slave.Close()
			go slave.Handle()
			<-doneChan
		}()
		master, mErr := NewMaster(listener, controlListener)
		assert.Nil(tt, mErr)
		defer master.Close()
		go func() {
			aConn, aErr := master.Accept()
			assert.Nil(tt, aErr)
			defer aConn.Close()
			<-doneChan
		}()
		aConn := network.Dial(listener.Addr().String())
		defer aConn.Close()
		doneChan <- struct{}{}
	})
}

func TestMaster_Dial(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		listener := network.ListenAny()
		defer listener.Close()
		controlListener := network.ListenAny()
		defer controlListener.Close()
		slaveConn := network.Dial(controlListener.Addr().String())
		defer slaveConn.Close()
		service := network.ListenAny()
		defer service.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		go func() {
			c, err := service.Accept()
			assert.Nil(tt, err)
			defer c.Close()
			_, err = c.Write([]byte(testMessage))
			assert.Nil(tt, err)
			<-doneChan
		}()
		go func() {
			slave, err := NewSlave(slaveConn)
			assert.Nil(tt, err)
			defer slave.Close()
			go slave.Handle()
			<-doneChan
		}()
		master, mErr := NewMaster(listener, controlListener)
		assert.Nil(tt, mErr)
		defer master.Close()
		serviceConn, dialErr := master.Dial("tcp", service.Addr().String())
		assert.Nil(tt, dialErr)
		defer serviceConn.Close()
		buffer := make([]byte, len(testMessage))
		_, rErr := serviceConn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, testMessage, string(buffer))
		doneChan <- struct{}{}
		doneChan <- struct{}{}
	})
}
