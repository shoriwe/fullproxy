package proxies

import (
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/shoriwe/fullproxy/v3/reverse"
	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func sshClient(t *testing.T) *ssh.Client {
	config := ssh.ClientConfig{
		User: "low",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, dErr := ssh.Dial("tcp", "127.0.0.1:2222", &config)
	assert.Nil(t, dErr)
	return client
}

func TestSSH_Addr(t *testing.T) {
	listener := network.ListenAny()
	client := sshClient(t)
	s := &SSH{
		Listener: listener,
		Client:   client,
	}
	defer s.Close()
	assert.NotNil(t, s.Addr())
}

func TestSSH_Serve(t *testing.T) {
	t.Run("Trigger KeepAlive", func(tt *testing.T) {
		listener := network.ListenAny()
		client := sshClient(t)
		s := &SSH{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: listener,
			Client:   client,
			Dial:     client.Dial,
		}
		defer s.Close()
		go s.Serve()
		time.Sleep(2 * time.Second)
		expect := httpexpect.Default(tt, "http://"+listener.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
	t.Run("Error KeepAlive", func(tt *testing.T) {
		listener := network.ListenAny()
		client := sshClient(t)
		s := &SSH{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: listener,
			Client:   client,
			Dial:     client.Dial,
		}

		go s.Serve()
		time.Sleep(2 * time.Second)
		s.Close()
		time.Sleep(2 * time.Second)
	})
}

func TestSSH_Serve_Local(t *testing.T) {
	t.Run("Basic", func(tt *testing.T) {
		listener := network.ListenAny()
		client := sshClient(t)
		s := &SSH{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: listener,
			Client:   client,
			Dial:     client.Dial,
		}
		defer s.Close()
		go s.Serve()
		expect := httpexpect.Default(tt, "http://"+listener.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
	t.Run("Reverse", func(tt *testing.T) {
		data := network.ListenAny()
		defer data.Close()
		control := network.ListenAny()
		defer control.Close()
		master := network.Dial(control.Addr().String())
		defer master.Close()
		listener := network.ListenAny()
		defer listener.Close()
		// Slave
		slave := &reverse.Slave{
			Listener: listener,
			Master:   master,
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
		client := sshClient(t)
		s := &SSH{
			Network:  "tcp",
			Address:  "echo:80",
			Listener: m,
			Client:   client,
			Dial:     client.Dial,
		}
		defer s.Close()
		go s.Serve()
		expect := httpexpect.Default(tt, "http://"+m.Addr().String())
		expect.GET("/").Expect().Status(http.StatusOK).Body().Contains("ECHO")
	})
}
