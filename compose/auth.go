package compose

import (
	"fmt"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

type Auth struct {
	Username   *string `yaml:"username,omitempty" json:"username,omitempty"`
	Password   *string `yaml:"password,omitempty" json:"password,omitempty"`
	PrivateKey *string `yaml:"privateKey,omitempty" json:"privateKey,omitempty"`
}

func (a *Auth) SSHClientConfig() (config *ssh.ClientConfig, err error) {
	if a.Username == nil {
		err = fmt.Errorf("no username provided")
	}
	if err == nil {
		config = new(ssh.ClientConfig)
		config.User = *a.Username
		if a.Password != nil {
			config.Auth = append(config.Auth, ssh.Password(*a.Password))
		}
		if a.PrivateKey != nil {
			// TODO: FIXME: This code doesn't work
			config.Auth = append(config.Auth, ssh.PublicKeys())
		}
	}
	return config, err
}

func (a *Auth) Socks5() (auth *proxy.Auth, err error) {
	if a.Username == nil {
		err = fmt.Errorf("no username provided")
	}
	if a.Password == nil {
		err = fmt.Errorf("no password provided")
	}
	if err == nil {
		auth = &proxy.Auth{
			User:     *a.Username,
			Password: *a.Password,
		}
	}
	return auth, err
}
