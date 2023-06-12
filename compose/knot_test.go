package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKnot_Compile_INVALID(t *testing.T) {
	knot := Knot{
		Type:    "INVALID",
		Network: "tcp",
		Address: "localhost:0",
	}
	_, err := knot.Compile()
	assert.NotNil(t, err)
}

func TestKnot_Compile_Forward(t *testing.T) {
	knot := Knot{
		Type:    KnotForward,
		Network: "tcp",
		Address: "localhost:0",
	}
	_, err := knot.Compile()
	assert.Nil(t, err)
}

func TestKnot_Compile_Socks5(t *testing.T) {
	t.Run("No Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSocks5,
			Network: "tcp",
			Address: "localhost:0",
		}
		_, err := knot.Compile()
		assert.Nil(t, err)
	})
	t.Run("Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSocks5,
			Network: "tcp",
			Address: "localhost:0",
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*knot.Auth.Username = "sulcud"
		*knot.Auth.Password = "password"
		_, err := knot.Compile()
		assert.Nil(t, err)
	})
	t.Run("Invalid Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSocks5,
			Network: "tcp",
			Address: "localhost:0",
			Auth: &Auth{
				Username: new(string),
			},
		}
		*knot.Auth.Username = "sulcud"
		_, err := knot.Compile()
		assert.NotNil(t, err)
	})
}

func TestKnot_Compile_SSH(t *testing.T) {
	t.Run("No Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSSH,
			Network: "tcp",
			Address: "localhost:0",
		}
		_, err := knot.Compile()
		assert.NotNil(t, err)
	})
	t.Run("Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSSH,
			Network: "tcp",
			Address: "localhost:0",
			Auth: &Auth{
				Username: new(string),
				Password: new(string),
			},
		}
		*knot.Auth.Username = "sulcud"
		*knot.Auth.Password = "password"
		_, err := knot.Compile()
		assert.Nil(t, err)
	})
	t.Run("Invalid Auth", func(tt *testing.T) {
		knot := Knot{
			Type:    KnotSSH,
			Network: "tcp",
			Address: "localhost:0",
			Auth:    &Auth{},
		}
		_, err := knot.Compile()
		assert.NotNil(t, err)
	})
}
