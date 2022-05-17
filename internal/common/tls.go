package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func SelfSignCertificate() ([]tls.Certificate, error) {
	var (
		priv *rsa.PrivateKey
		cert []byte
		err  error
	)
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(now.Unix()),
		Subject:               pkix.Name{},
		NotBefore:             now,
		NotAfter:              now.AddDate(1000, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	cert, err = x509.CreateCertificate(rand.Reader, template, template,
		priv.Public(), priv)
	if err != nil {
		return nil, err
	}

	var outCert tls.Certificate
	outCert.Certificate = append(outCert.Certificate, cert)
	outCert.PrivateKey = priv
	return []tls.Certificate{outCert}, nil

}
