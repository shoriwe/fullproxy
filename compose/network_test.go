package compose

import (
	"net"
	"testing"

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
		*l.Address = "localhost:2222"
		*l.Data.Network = "tcp"
		*l.Data.Address = "localhost:0"
		*l.Auth.Username = "low"
		*l.Auth.Password = "password"
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
		*l.Address = "localhost:2222"
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
		*l.Auth.Username = "low"
		*l.Auth.Password = "password"
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
