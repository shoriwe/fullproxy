package compose

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

type Auth struct {
	Username   *string `yaml:"username,omitempty" json:"username,omitempty"`
	Password   *string `yaml:"password,omitempty" json:"password,omitempty"`
	PrivateKey *string `yaml:"privateKey,omitempty" json:"privateKey,omitempty"`
	ServerKey  *string `yaml:"serverKey,omitempty" json:"serverKey,omitempty"`
}

func (a *Auth) getHostKeyCallback() (ssh.HostKeyCallback, error) {
	if a.ServerKey == nil {
		return ssh.InsecureIgnoreHostKey(), nil
	}
	// TODO: Test this code
	serverKey, err := ioutil.ReadFile(*a.ServerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key file: %v", err)
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey(serverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server key: %v", err)
	}
	return ssh.FixedHostKey(hostKey), nil
}

func (a *Auth) SSHClientConfig() (config *ssh.ClientConfig, err error) {
	if a.Username == nil {
		err = fmt.Errorf("no username provided")
	}
	if err == nil {
		config = new(ssh.ClientConfig)
		config.HostKeyCallback, err = a.getHostKeyCallback()
		if err == nil {
			config.User = *a.Username
			if a.Password != nil {
				config.Auth = append(config.Auth, ssh.Password(*a.Password))
			}
			// TODO: Test this code
			if a.PrivateKey != nil {
				key, err := ioutil.ReadFile(*a.PrivateKey)
				if err != nil {
					return nil, fmt.Errorf("failed to read private key file: %v", err)
				}
				signer, err := ssh.ParsePrivateKey(key)
				if err != nil {
					return nil, fmt.Errorf("failed to parse private key: %v", err)
				}
				config.Auth = append(config.Auth, ssh.PublicKeys(signer))
			}
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
