package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse_Master(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		r := Reverse{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
			Controller: Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
		}
		m, err := r.Master()
		assert.Nil(tt, err)
		defer m.Close()
	})
	t.Run("Invalid User listener", func(tt *testing.T) {
		r := Reverse{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:99999999",
			},
			Controller: Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
		}
		_, err := r.Master()
		assert.NotNil(tt, err)
	})
	t.Run("Invalid Controller listener", func(tt *testing.T) {
		r := Reverse{
			Listener: &Listener{
				Network: "tcp",
				Address: "localhost:0",
			},
			Controller: Listener{
				Network: "tcp",
				Address: "localhost:9999999999",
			},
		}
		_, err := r.Master()
		assert.NotNil(tt, err)
	})
}
