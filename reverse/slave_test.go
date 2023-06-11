package reverse

import (
	"encoding/gob"
	"net"
	"testing"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestSlave_init(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		go func() {
			s := &Slave{Master: master, Dial: net.Dial}
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
			s := &Slave{Master: master, Dial: net.Dial}
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

func TestSlave_HandleAccept(t *testing.T) {
	t.Run("Succeed", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		sListener := network.ListenAny()
		defer sListener.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		testMessage := []byte("TEST")
		go func() {
			m := &Master{
				Data:    data,
				Control: control,
			}
			defer m.Close()
			conn, cErr := m.SlaveAccept()
			assert.Nil(tt, cErr)
			defer conn.Close()
			_, wErr := conn.Write(testMessage)
			assert.Nil(tt, wErr)
			<-doneChan
		}()
		s := &Slave{Master: master, Listener: sListener, Dial: net.Dial}
		defer s.Close()
		assert.Nil(tt, s.init())
		go func() {
			job, aErr := s.Control.Accept()
			assert.Nil(tt, aErr)
			defer job.Close()
			var req Request
			dErr := gob.NewDecoder(job).Decode(&req)
			assert.Nil(tt, dErr)
			go s.HandleAccept(job, &req)
			<-doneChan
		}()
		conn, dErr := net.Dial(sListener.Addr().Network(), sListener.Addr().String())
		assert.Nil(tt, dErr)
		buffer := make([]byte, len(testMessage))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, testMessage, buffer)
		doneChan <- struct{}{}
		doneChan <- struct{}{}
	})
	t.Run("No listener", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		doneChan := make(chan struct{}, 1)
		defer close(doneChan)
		go func() {
			m := &Master{
				Data:    data,
				Control: control,
			}
			defer m.Close()
			_, cErr := m.SlaveAccept()
			assert.NotNil(tt, cErr)
			<-doneChan
		}()
		s := &Slave{Master: master, Dial: net.Dial}
		defer s.Close()
		assert.Nil(tt, s.init())
		job, aErr := s.Control.Accept()
		assert.Nil(tt, aErr)
		defer job.Close()
		var req Request
		dErr := gob.NewDecoder(job).Decode(&req)
		assert.Nil(tt, dErr)
		assert.NotNil(tt, s.HandleAccept(job, &req))
		doneChan <- struct{}{}
	})
}

func TestSlave_Handle(t *testing.T) {
	t.Run("Accept", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		sListener := network.ListenAny()
		defer sListener.Close()
		doneChan := make(chan struct{}, 2)
		defer close(doneChan)
		testMessage := []byte("TEST")
		go func() {
			m := &Master{
				Data:    data,
				Control: control,
			}
			defer m.Close()
			conn, cErr := m.SlaveAccept()
			assert.Nil(tt, cErr)
			defer conn.Close()
			_, wErr := conn.Write(testMessage)
			assert.Nil(tt, wErr)
			<-doneChan
		}()
		s := &Slave{Master: master, Listener: sListener, Dial: net.Dial}
		defer s.Close()
		assert.Nil(tt, s.init())
		go func() {
			job, aErr := s.Control.Accept()
			assert.Nil(tt, aErr)
			defer job.Close()
			go s.Handle(job)
			<-doneChan
		}()
		conn, dErr := net.Dial(sListener.Addr().Network(), sListener.Addr().String())
		assert.Nil(tt, dErr)
		buffer := make([]byte, len(testMessage))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, testMessage, buffer)
		doneChan <- struct{}{}
		doneChan <- struct{}{}
	})
	t.Run("Dial", func(tt *testing.T) {
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
			m := &Master{
				Data:    data,
				Control: control,
			}
			defer m.Close()
			conn, cErr := m.SlaveDial(service.Addr().Network(), service.Addr().String())
			assert.Nil(tt, cErr)
			defer conn.Close()
			_, wErr := conn.Write(testMessage)
			assert.Nil(tt, wErr)
			<-doneChan
		}()
		s := &Slave{Master: master, Dial: net.Dial}
		defer s.Close()
		assert.Nil(tt, s.init())
		go func() {
			job, aErr := s.Control.Accept()
			assert.Nil(tt, aErr)
			defer job.Close()
			go s.Handle(job)
			<-doneChan
		}()
		conn, dErr := service.Accept()
		assert.Nil(tt, dErr)
		buffer := make([]byte, len(testMessage))
		_, rErr := conn.Read(buffer)
		assert.Nil(tt, rErr)
		assert.Equal(tt, testMessage, buffer)
		doneChan <- struct{}{}
		doneChan <- struct{}{}
	})
	t.Run("Invalid action", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		sListener := network.ListenAny()
		defer sListener.Close()
		doneChan := make(chan struct{}, 1)
		defer close(doneChan)
		go func() {
			m := &Master{
				Data:    data,
				Control: control,
			}
			defer m.Close()
			assert.Nil(tt, m.init())
			conn, cErr := m.cSession.Open()
			assert.Nil(tt, cErr)
			defer conn.Close()
			assert.Nil(tt, gob.NewEncoder(conn).Encode(Request{
				Action: Invalid,
			}))
			<-doneChan
		}()
		s := &Slave{Master: master, Listener: sListener, Dial: net.Dial}
		defer s.Close()
		assert.Nil(tt, s.init())
		job, aErr := s.Control.Accept()
		assert.Nil(tt, aErr)
		defer job.Close()
		assert.NotNil(tt, s.Handle(job))
		doneChan <- struct{}{}
	})
}
