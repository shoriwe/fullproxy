package compose

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListener_setupBasic(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		ll, err := l.setupBasic(net.Listen)
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Network", func(tt *testing.T) {
		l := Listener{
			Address: new(string),
		}
		*l.Address = "localhost:0"
		_, err := l.setupBasic(net.Listen)
		assert.NotNil(tt, err)
	})
	t.Run("No Address", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
		}
		*l.Network = "tcp"
		_, err := l.setupBasic(net.Listen)
		assert.NotNil(tt, err)
	})
}

func TestListener_setupMaster(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Listener{
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		ll, err := l.setupMaster()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Data", func(tt *testing.T) {
		l := Listener{
			Control: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		_, err := l.setupMaster()
		assert.NotNil(tt, err)
	})
	t.Run("No Control", func(tt *testing.T) {
		l := Listener{
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupMaster()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Data", func(tt *testing.T) {
		l := Listener{
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Control.Network = "tcp"
		*l.Control.Address = "localhost:0"
		_, err := l.setupMaster()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Control", func(tt *testing.T) {
		l := Listener{
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupMaster()
		assert.NotNil(tt, err)
	})
}

func TestListener_setupSSH(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = "localhost:2222"
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = "low"
		*l.Auth.Password = "password"
		ll, err := l.setupSSH()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("No Network", func(tt *testing.T) {
		l := Listener{}
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
	t.Run("No Address", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
		}
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
	t.Run("No Data", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
		}
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
	t.Run("No Auth", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
		}
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Auth", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{},
		}
		*l.Network = "tcp"
		*l.Address = "localhost:2222"
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
	t.Run("Dial denied", func(tt *testing.T) {
		l := Listener{
			Network: new(string),
			Address: new(string),
			Data: &Listener{
				Type:    ListenerBasic,
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
		*l.Auth.Username = "low"
		*l.Auth.Password = "password"
		_, err := l.setupSSH()
		assert.NotNil(tt, err)
	})
}

func TestListener_Listen(t *testing.T) {
	t.Run(ListenerBasic, func(tt *testing.T) {
		l := Listener{
			Type:    ListenerBasic,
			Network: new(string),
			Address: new(string),
		}
		*l.Network = "tcp"
		*l.Address = "localhost:0"
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run(ListenerMaster, func(tt *testing.T) {
		l := Listener{
			Type: ListenerMaster,
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Control: &Listener{
				Type:    ListenerBasic,
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
	t.Run(ListenerMaster, func(tt *testing.T) {
		l := Listener{
			Type:    ListenerSSH,
			Network: new(string),
			Address: new(string),
			Data: &Listener{
				Type:    ListenerBasic,
				Network: new(string),
				Address: new(string),
			},
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*l.Network = "tcp"
		*l.Address = "localhost:2222"
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = "low"
		*l.Auth.Password = "password"
		ll, err := l.Listen()
		assert.Nil(tt, err)
		defer ll.Close()
	})
	t.Run("UNKNOWN", func(tt *testing.T) {
		l := Listener{
			Type: "UNKNOWN",
		}
		_, err := l.Listen()
		assert.NotNil(tt, err)
	})
	t.Run("Self signed TLS", func(tt *testing.T) {
		l := Listener{
			Type:    ListenerBasic,
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
