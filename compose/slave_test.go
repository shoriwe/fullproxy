package compose

import (
	"testing"

	"github.com/shoriwe/fullproxy/v3/reverse"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestSlave_setupSlave(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		m := reverse.Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		go m.Accept()
		s := Slave{
			MasterNetwork: control.Addr().Network(),
			MasterAddress: control.Addr().String(),
			MasterDialer: Network{
				Type: NetworkBasic,
			},
			Dialer: Network{
				Type: NetworkBasic,
			},
		}
		slave, err := s.setupSlave()
		assert.Nil(tt, err)
		go slave.Serve()
		conn := network.Dial(data.Addr().String())
		defer conn.Close()
	})
	t.Run("With Listener", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		s := Slave{
			MasterNetwork: control.Addr().Network(),
			MasterAddress: control.Addr().String(),
			MasterDialer: Network{
				Type: NetworkBasic,
			},
			Dialer: Network{
				Type: NetworkBasic,
			},
			Listener: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*s.Listener.Network = "tcp"
		*s.Listener.Address = "localhost:0"
		// Slave
		slave, err := s.setupSlave()
		assert.Nil(tt, err)
		go slave.Serve()
		go func() {
			conn := network.Dial(s.Listener.listener.Addr().String())
			defer conn.Close()
		}()
		// Master
		m := reverse.Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		conn, err := m.SlaveAccept()
		assert.Nil(tt, err)
		defer conn.Close()
	})
}

func TestSlave_Serve(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		m := reverse.Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		go m.Accept()
		s := Slave{
			MasterNetwork: control.Addr().Network(),
			MasterAddress: control.Addr().String(),
			MasterDialer: Network{
				Type: NetworkBasic,
			},
			Dialer: Network{
				Type: NetworkBasic,
			},
		}
		go s.Serve()
		conn := network.Dial(data.Addr().String())
		defer conn.Close()
	})
}
