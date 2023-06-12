package proxies

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v4/reverse"
	"github.com/shoriwe/fullproxy/v4/sshd"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestSSH_Addr(t *testing.T) {
	listener := network.ListenAny()
	defer listener.Close()
	client := sshd.DefaultClient(t)
	defer client.Close()
	go sshd.KeepAlive(client)
	s := &Forward{
		Listener: listener,
	}
	defer s.Close()
	assert.NotNil(t, s.Addr())
}

func TestSSH_Serve(t *testing.T) {
	t.Run("Trigger KeepAlive", func(tt *testing.T) {
		listener := network.ListenAny()
		defer listener.Close()
		client := sshd.DefaultClient(tt)
		defer client.Close()
		go sshd.KeepAlive(client)
		s := &Forward{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: listener,
			Dial:     client.Dial,
		}
		defer s.Close()
		go s.Serve()
		expect := httpexpect.Default(tt, "http://"+listener.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}

func TestSSH_Serve_Local(t *testing.T) {
	t.Run("Basic", func(tt *testing.T) {
		listener := network.ListenAny()
		sshClient := sshd.DefaultClient(tt)
		go sshd.KeepAlive(sshClient)
		s := &Forward{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: listener,
			Dial:     sshClient.Dial,
		}
		defer s.Close()
		go s.Serve()
		expect := httpexpect.Default(tt, "http://"+listener.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
	t.Run("Reverse", func(tt *testing.T) {
		// Master
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		// Slave
		masterConn := network.Dial(control.Addr().String())
		defer masterConn.Close()
		slaveListener := network.ListenAny()
		defer slaveListener.Close()
		// SSH Conn
		sshClient := sshd.DefaultClient(tt)
		go sshd.KeepAlive(sshClient)
		// Slave
		slave := &reverse.Slave{
			Listener: slaveListener,
			Master:   masterConn,
		}
		defer slave.Close()
		go slave.Serve()
		// Master
		m := &reverse.Master{
			Data:    data,
			Control: control,
		}
		defer m.Close()
		//
		s := &Forward{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: m,
			Dial:     sshClient.Dial,
		}
		defer s.Close()
		go s.Serve()
		expect := httpexpect.Default(tt, "http://"+m.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}
