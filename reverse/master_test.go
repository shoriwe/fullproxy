package reverse

import (
	"net"
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestMaster_init(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		go func() {
			s := &Slave{Master: master}
			defer s.Close()
			assert.Nil(tt, s.init())
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		assert.Nil(tt, m.init())
	})
	t.Run("Twice", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		go func() {
			master := network.Dial(control.Addr().String())
			defer master.Close()
			s := &Slave{Master: master}
			defer s.Close()
			assert.Nil(tt, s.init())
			assert.Nil(tt, s.init())
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		assert.Nil(tt, m.init())
		assert.Nil(tt, m.init())
	})
}

func TestMaster_Addr(t *testing.T) {
	data := network.ListenAny()
	defer data.Close()
	control := network.ListenAny()
	defer control.Close()
	m := &Master{
		Data:    data,
		Control: control,
	}
	assert.NotNil(t, m.Addr())
}
func TestMaster_Accept(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		go func() {
			s := &Slave{Master: master}
			defer s.Close()
			go s.Serve()
			<-doneChan
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		go func() {
			aConn, aErr := m.Accept()
			assert.Nil(tt, aErr)
			defer aConn.Close()
			<-doneChan
		}()
		aConn := network.Dial(data.Addr().String())
		defer aConn.Close()
		doneChan <- struct{}{}
	})
}

func TestMaster_SlaveDial(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		service := network.ListenAny()
		defer service.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		testMessage := []byte("TEST")
		go func() {
			c, err := service.Accept()
			assert.Nil(tt, err)
			defer c.Close()
			_, err = c.Write(testMessage)
			assert.Nil(tt, err)
			<-doneChan
		}()
		go func() {
			s := &Slave{Master: master}
			defer s.Close()
			go s.Serve()
			<-doneChan
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		serviceConn, dialErr := m.SlaveDial("tcp", service.Addr().String())
		assert.Nil(tt, dialErr)
		defer serviceConn.Close()
		buffer := make([]byte, len(testMessage))
		_, rErr := serviceConn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, testMessage, buffer)
		doneChan <- struct{}{}
		doneChan <- struct{}{}
	})
	t.Run("Not listening", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		service := network.ListenAny()
		assert.Nil(tt, service.Close())
		doneChan := make(chan struct{}, 1)
		defer close(doneChan)
		go func() {
			s := &Slave{Master: master}
			defer s.Close()
			go s.Serve()
			<-doneChan
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		_, dialErr := m.SlaveDial("tcp", service.Addr().String())
		assert.NotNil(tt, dialErr)
		doneChan <- struct{}{}
	})
}

func TestMaster_SlaveAccept(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		sListener := network.ListenAny() // Slave listener
		defer sListener.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		go func() {
			s := &Slave{Master: master, Listener: sListener}
			defer s.Close()
			go s.Serve()
			<-doneChan
		}()
		m := &Master{
			Data:    data,
			Control: control,
		}
		go func() {
			conn, dErr := net.Dial(sListener.Addr().Network(), sListener.Addr().String())
			assert.Nil(tt, dErr)
			defer conn.Close()
		}()
		defer m.Close()
		client, dialErr := m.SlaveAccept()
		assert.Nil(tt, dialErr)
		defer client.Close()
	})
}
