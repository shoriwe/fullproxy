package reverse

import (
	"net"
	"sync"
	"testing"

	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestMaster_init(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		masterConn := network.Dial(control.Addr().String())
		defer masterConn.Close()
		go func() {
			slave := &Slave{Master: masterConn, Dial: net.Dial}
			defer slave.Close()
			assert.Nil(tt, slave.init())
		}()
		master := &Master{
			Data:    data,
			Control: control,
		}
		defer master.Close()
		assert.Nil(tt, master.init())
	})
	t.Run("Twice", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		masterConn := network.Dial(control.Addr().String())
		defer masterConn.Close()
		go func() {
			slave := &Slave{Master: masterConn, Dial: net.Dial}
			defer slave.Close()
			assert.Nil(tt, slave.init())
			assert.Nil(tt, slave.init())
		}()
		master := &Master{
			Data:    data,
			Control: control,
		}
		defer master.Close()
		assert.Nil(tt, master.init())
		assert.Nil(tt, master.init())
	})
}

func TestMaster_Addr(t *testing.T) {
	data := network.ListenAny()
	defer data.Close()
	control := network.ListenAny()
	defer control.Close()
	master := &Master{
		Data:    data,
		Control: control,
	}
	defer master.Close()
	assert.NotNil(t, master.Addr())
}
func TestMaster_Accept(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		masterConn := network.Dial(control.Addr().String())
		defer masterConn.Close()
		// Slave
		slave := &Slave{Master: masterConn, Dial: net.Dial}
		defer slave.Close()
		go slave.Serve()
		// Master
		master := &Master{
			Data:    data,
			Control: control,
		}
		defer master.Close()
		// Producer
		testMsg := []byte("MSG")
		var wg sync.WaitGroup
		defer wg.Wait()
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := master.Accept()
			assert.Nil(tt, err)
			defer conn.Close()
			// Write
			_, err = conn.Write(testMsg)
			assert.Nil(tt, err)
			// Read
			buffer := make([]byte, len(testMsg))
			_, err = conn.Read(buffer)
			assert.Nil(tt, err)
		}()
		// Consumer
		conn := network.Dial(data.Addr().String())
		defer conn.Close()
		// Read
		buffer := make([]byte, len(testMsg))
		_, err := conn.Read(buffer)
		assert.Nil(tt, err)
		// Write
		_, err = conn.Write(testMsg)
		assert.Nil(tt, err)
	})
}

func TestMaster_SlaveDial(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		masterConn := network.Dial(control.Addr().String())
		defer masterConn.Close()
		service := network.ListenAny()
		defer service.Close()
		// Slave
		slave := &Slave{Master: masterConn, Dial: net.Dial}
		defer slave.Close()
		go slave.Serve()
		// Master
		master := &Master{
			Data:    data,
			Control: control,
		}
		defer master.Close()
		// Producer
		testMessage := []byte("TEST")
		var wg sync.WaitGroup
		defer wg.Wait()
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := service.Accept()
			assert.Nil(tt, err)
			defer conn.Close()
			// Write
			_, err = conn.Write(testMessage)
			assert.Nil(tt, err)
			// Read
			buffer := make([]byte, len(testMessage))
			_, err = conn.Read(buffer)
			assert.Nil(tt, err)
		}()
		// Consumer
		conn, err := master.SlaveDial("tcp", service.Addr().String())
		assert.Nil(tt, err)
		defer conn.Close()
		// Read
		buffer := make([]byte, len(testMessage))
		_, err = conn.Read(buffer)
		assert.Nil(tt, err)
		assert.Equal(tt, testMessage, buffer)
		// Write
		_, err = conn.Write(testMessage)
		assert.Nil(tt, err)
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
			s := &Slave{Master: master, Dial: net.Dial}
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
			s := &Slave{Master: master, Listener: sListener, Dial: net.Dial}
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
