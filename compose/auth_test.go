package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_SSHClientConfig(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		auth := Auth{
			Username: new(string),
			Password: new(string),
		}
		*auth.Username = "sulcud"
		*auth.Password = "password"
		_, err := auth.SSHClientConfig()
		assert.Nil(tt, err)
	})
	t.Run("Private Key", func(tt *testing.T) {
		// TODO: Implement me!
	})
	t.Run("No Username", func(tt *testing.T) {
		auth := Auth{}
		_, err := auth.SSHClientConfig()
		assert.NotNil(tt, err)
	})
}

func TestAuth_Socks5(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		auth := Auth{
			Username: new(string),
			Password: new(string),
		}
		*auth.Username = "sulcud"
		*auth.Password = "password"
		_, err := auth.Socks5()
		assert.Nil(tt, err)
	})
	t.Run("No Username", func(tt *testing.T) {
		auth := Auth{}
		_, err := auth.Socks5()
		assert.NotNil(tt, err)
	})
	t.Run("No Password", func(tt *testing.T) {
		auth := Auth{
			Username: new(string),
		}
		*auth.Username = "sulcud"
		_, err := auth.Socks5()
		assert.NotNil(tt, err)
	})
}
