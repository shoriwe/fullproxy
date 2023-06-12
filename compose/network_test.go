package compose

import (
	"net"
	"testing"

	"github.com/shoriwe/fullproxy/v4/reverse"
	"github.com/shoriwe/fullproxy/v4/sshd"
	"github.com/shoriwe/fullproxy/v4/utils/network"
	"github.com/stretchr/testify/assert"
)

func TestNetwork_setupBasicListener(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		ll, err := l.setupBasicListener(net.Listen)
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Network", func(tt *testing.T) {
		l := Network{
			Address: new(string),
		}
		*l.Address = "localhost:0"
		_, err := l.setupBasicListener(net.Listen)
		assert.NotNil(tt, err)
	})
	t.Run("No Address", func(tt *testing.T) {
		l := Network{
			Network: new(string),
		}
		*l.Network = "tcp"
		_, err := l.setupBasicListener(net.Listen)
		assert.NotNil(tt, err)
	})
}

func TestNetwork_setupMasterListener(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Network{
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		ll, err := l.setupMasterListener()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Data", func(tt *testing.T) {
		l := Network{
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		_, err := l.setupMasterListener()
		assert.NotNil(tt, err)
	})
	t.Run("No Control", func(tt *testing.T) {
		l := Network{
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupMasterListener()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Data", func(tt *testing.T) {
		l := Network{
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		_, err := l.setupMasterListener()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Control", func(tt *testing.T) {
		l := Network{
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupMasterListener()
		assert.NotNil(tt, err)
	})
	t.Run("SlaveListener", func(tt *testing.T) {
		l := Network{
			SlaveListener: new(bool),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.SlaveListener = true
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		// Request the slave listener
		ll, err := l.setupMasterListener()
		assert.Nil(tt, err)
		defer ll.Close()
		//
		checkCh := make(chan struct{}, 1)
		ssL := network.ListenAny()
		defer ssL.Close()
		// Start slave
		go func() {
			masterConn := network.Dial(l.master.Control.Addr().String())
			s := &reverse.Slave{
				Listener: ssL,
				Dial:     net.Dial,
				Master:   masterConn,
			}
			go s.Serve()
			<-checkCh
		}()
		testMsg := []byte("TEST")
		// Start Dial to slave listener
		go func() {
			conn := network.Dial(ssL.Addr().String())
			defer conn.Close()
			_, err := conn.Write(testMsg)
			assert.Nil(tt, err)
		}()
		// Accept connection
		conn, err := ll.Accept()
		assert.Nil(tt, err)
		defer conn.Close()
		buffer := make([]byte, len(testMsg))
		_, err = conn.Read(buffer)
		assert.Nil(tt, err)
		assert.Equal(tt, testMsg, buffer)
		checkCh <- struct{}{}
	})
}

func TestNetwork_setupSSHListener(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = sshd.DefaultAddr
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = sshd.SSHDefaultUser
		*l.Auth.Password = sshd.SSHDefaultPassword
		ll, err := l.setupSSHListener()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Network", func(tt *testing.T) {
		l := Network{}
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
	t.Run("No Address", func(tt *testing.T) {
		l := Network{
			Network: new(string),
		}
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
	t.Run("No Data", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
		}
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
	t.Run("No Auth", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Auth", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{},
		}
		*l.Network = "tcp"
		*l.Address = sshd.DefaultAddr
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
	t.Run("Dial denied", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = "1111111111111111111111111111111"
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = sshd.SSHDefaultUser
		*l.Auth.Password = sshd.SSHDefaultPassword
		_, err := l.setupSSHListener()
		assert.NotNil(tt, err)
	})
}

func TestNetwork_Listen(t *testing.T) {
	t.Run(NetworkBasic, func(tt *testing.T) {
		l := Network{
			Type:    NetworkBasic,
			Network: new(string),
			Address: new(string),
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run(NetworkMaster, func(tt *testing.T) {
		l := Network{
			Type: NetworkMaster,
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run(NetworkMaster, func(tt *testing.T) {
		l := Network{
			Type:    NetworkSSH,
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = sshd.DefaultAddr
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = sshd.SSHDefaultUser
		*l.Auth.Password = sshd.SSHDefaultPassword
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("UNKNOWN", func(tt *testing.T) {
		l := Network{
			Type: "UNKNOWN",
		}
		_, err := l.Listen()
		assert.NotNil(tt, err)
	})
	t.Run("Self signed TLS", func(tt *testing.T) {
		l := Network{
			Type:    NetworkBasic,
			Network: new(string),
			Address: new(string),
			Crypto: &Crypto{
				Mode:       CryptoTLS,
				SelfSigned: true,
			},
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
}

func TestNetwork_setupMasterDialFunc(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Network{
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		dialFunc, err := l.setupMasterDialFunc()
		assert.Nil(tt, err)
		master, err := l.getMaster()
		assert.Nil(tt, err)
		defer master.Close()
		assert.NotNil(tt, dialFunc)
	})
}

func TestNetwork_setupSSHDialFunc(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Network{
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = sshd.DefaultAddr
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = sshd.SSHDefaultUser
		*l.Auth.Password = sshd.SSHDefaultPassword
		dialFunc, err := l.setupSSHDialFunc()
		assert.Nil(tt, err)
		ssh, err := l.getSSHConn()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
		go ssh.Close()
	})
}

func TestNetwork_DialFunc(t *testing.T) {
	t.Run(NetworkBasic, func(tt *testing.T) {
		l := Network{
			Type:    NetworkBasic,
			Network: new(string),
			Address: new(string),
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		dialFunc, err := l.DialFunc()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
	})
	t.Run(NetworkMaster, func(tt *testing.T) {
		l := Network{
			Type: NetworkMaster,
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		dialFunc, err := l.DialFunc()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
		master, err := l.getMaster()
		assert.Nil(tt, err)
		defer master.Close()
	})
	t.Run(NetworkMaster, func(tt *testing.T) {
		l := Network{
			Type:    NetworkSSH,
			Network: new(string),
			Address: new(string),
			Data: &Network{
				Type:    NetworkBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = sshd.DefaultAddr
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = sshd.SSHDefaultUser
		*l.Auth.Password = sshd.SSHDefaultPassword
		dialFunc, err := l.DialFunc()
		assert.Nil(tt, err)
		assert.NotNil(tt, dialFunc)
		master, err := l.getSSHConn()
		assert.Nil(tt, err)
		defer master.Close()
	})
	t.Run("UNKNOWN", func(tt *testing.T) {
		l := Network{
			Type: "UNKNOWN",
		}
		_, err := l.DialFunc()
		assert.NotNil(tt, err)
	})
}
