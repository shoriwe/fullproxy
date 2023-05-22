package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/shoriwe/fullproxy/v3/utils/network"
	"github.com/stretchr/testify/assert"
)

func testTLS(t *testing.T, serverConfig, clientConfig *tls.Config) {
	l := network.ListenAny()
	defer l.Close()
	ltls := tls.NewListener(l, serverConfig)
	doneCh := make(chan struct{}, 1)
	msg := []byte("TEST")
	go func() {
		conn, dErr := tls.Dial(ltls.Addr().Network(), ltls.Addr().String(), clientConfig)
		assert.Nil(t, dErr)
		defer conn.Close()
		_, wErr := conn.Write(msg)
		assert.Nil(t, wErr)
		<-doneCh
	}()
	conn, aErr := ltls.Accept()
	assert.Nil(t, aErr)
	defer conn.Close()
	buffer := make([]byte, len(msg))
	_, rErr := conn.Read(buffer)
	assert.Nil(t, rErr)
	assert.Equal(t, msg, buffer)
	doneCh <- struct{}{}
}

func TestTLS(t *testing.T) {
	t.Run("In memory", func(tt *testing.T) {
		c := TLS{}
		config, cErr := c.Config()
		assert.Nil(tt, cErr)
		// -- The real test is made here
		testTLS(t, config, config)
	})
	t.Run("InsecureSkipVerify", func(tt *testing.T) {
		c1 := TLS{}
		serverConfig, cErr := c1.Config()
		assert.Nil(tt, cErr)
		c2 := TLS{
			InsecureSkipVerify: new(bool),
		}
		*c2.InsecureSkipVerify = true
		clientConfig, cErr := c2.Config()
		assert.Nil(tt, cErr)
		// -- The real test is made here
		testTLS(t, serverConfig, clientConfig)
	})
	t.Run("Cert and Key", func(tt *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Panicf("Failed to generate private key: %v", err)
		}
		sN, sNErr := rand.Int(rand.Reader, big.NewInt(0xFF_FF_FF_FF_FF_FF_FF))
		if sNErr != nil {
			log.Panicf("invalid serial number %v", sNErr)
		}
		template := &x509.Certificate{
			SerialNumber:          sN,
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(120, 0, 0),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}
		certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
		if err != nil {
			log.Panicf("Failed to create certificate: %v", err)
		}
		privateKeyPEM := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		}
		certPEM := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certDER,
		}
		// -- Key
		key, tErr := os.CreateTemp("", "*")
		assert.Nil(tt, tErr)
		defer os.RemoveAll(key.Name())
		eErr := pem.Encode(key, privateKeyPEM)
		assert.Nil(tt, eErr)
		key.Close()
		// -- Cert
		cert, tErr := os.CreateTemp("", "*")
		assert.Nil(tt, tErr)
		defer os.RemoveAll(cert.Name())
		eErr = pem.Encode(cert, certPEM)
		assert.Nil(tt, eErr)
		cert.Close()
		// -- File locations
		keyName := key.Name()
		certName := cert.Name()
		c1 := TLS{
			KeyFile:  &keyName,
			CertFile: &certName,
		}
		serverConfig, cErr := c1.Config()
		assert.Nil(tt, cErr)
		c2 := TLS{
			InsecureSkipVerify: new(bool),
		}
		*c2.InsecureSkipVerify = true
		clientConfig, cErr := c2.Config()
		assert.Nil(tt, cErr)
		// -- The real test is made here
		testTLS(t, serverConfig, clientConfig)
	})
}
