package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"
)

func DefaultTLSConfig() *tls.Config {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	sN, _ := rand.Int(rand.Reader, big.NewInt(0xFF_FF_FF_FF_FF_FF_FF))
	template := &x509.Certificate{
		SerialNumber:          sN,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(120, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}
	privateKeyPEMBytes := pem.EncodeToMemory(privateKeyPEM)
	certPEMBytes := pem.EncodeToMemory(certPEM)
	tlsCert, _ := tls.X509KeyPair(certPEMBytes, privateKeyPEMBytes)
	return &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: true,
	}
}

func TempCertKey() (string, string) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	sN, _ := rand.Int(rand.Reader, big.NewInt(0xFF_FF_FF_FF_FF_FF_FF))
	template := &x509.Certificate{
		SerialNumber:          sN,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(120, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}
	privateKeyFile, _ := ioutil.TempFile("", "privatekey*.pem")
	defer privateKeyFile.Close()
	pem.Encode(privateKeyFile, privateKeyPEM)
	certFile, _ := ioutil.TempFile("", "certificate*.pem")
	defer certFile.Close()
	pem.Encode(certFile, certPEM)
	return certFile.Name(), privateKeyFile.Name()
}
