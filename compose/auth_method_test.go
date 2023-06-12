package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMethod_Socks5(t *testing.T) {
	t.Run("Raw", func(tt *testing.T) {
		a := AuthMethod{
			Raw: map[string]string{
				"sulcud": "password",
			},
		}
		auth, err := a.Socks5()
		assert.Nil(tt, err)
		assert.NotNil(tt, auth)
	})
	t.Run("Nil", func(tt *testing.T) {
		a := AuthMethod{}
		auth, err := a.Socks5()
		assert.Nil(tt, err)
		assert.Nil(tt, auth)
	})
}

func TestAuthMethods_Socks5(t *testing.T) {
	t.Run("Raw", func(tt *testing.T) {
		a := AuthMethod{
			Raw: map[string]string{
				"sulcud": "password",
			},
		}
		as := AuthMethods{a}
		authMethods, err := as.Socks5()
		assert.Nil(tt, err)
		assert.NotNil(tt, authMethods)
	})
	t.Run("Nil", func(tt *testing.T) {
		as := AuthMethods{}
		authMethods, err := as.Socks5()
		assert.Nil(tt, err)
		assert.Nil(tt, authMethods)
	})
}
