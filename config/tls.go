package config

import (
	"crypto/tls"

	fullproxyCrypto "github.com/shoriwe/fullproxy/v3/utils/crypto"
)

type TLS struct {
	InsecureSkipVerify *bool   `yaml:"ignore"`
	CertFile           *string `yaml:"cert"`
	KeyFile            *string `yaml:"key"`
}

func (t *TLS) Config() (*tls.Config, error) {
	if t.InsecureSkipVerify == nil && t.CertFile == nil && t.KeyFile == nil {
		return fullproxyCrypto.DefaultTLSConfig(), nil
	}
	config := &tls.Config{}
	if t.InsecureSkipVerify != nil {
		config.InsecureSkipVerify = *t.InsecureSkipVerify
	}
	if t.CertFile != nil && t.KeyFile != nil {
		cert, lErr := tls.LoadX509KeyPair(*t.CertFile, *t.KeyFile)
		if lErr != nil {
			return nil, lErr
		}
		config.Certificates = append(config.Certificates, cert)
	}
	return config, nil
}
