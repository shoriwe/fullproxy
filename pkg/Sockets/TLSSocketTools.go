package Sockets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"time"
)

func CreateMasterSlaveCertificateTemplate() *x509.Certificate {
	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"FullProxy"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(1, 0, 0),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
}

func CreateCertificate(template *x509.Certificate) ([]byte, *rsa.PrivateKey, error) {
	privateKey, keyGenerationError := rsa.GenerateKey(rand.Reader, 2048)
	if keyGenerationError == nil {
		certificate, certificateCreationError := x509.CreateCertificate(
			rand.Reader,
			template,
			template,
			&privateKey.PublicKey,
			privateKey)
		if certificateCreationError == nil {
			return certificate, privateKey, nil
		}
		return nil, nil, certificateCreationError
	}
	return nil, nil, keyGenerationError

}

func CreateTLSConfiguration(template *x509.Certificate) (*tls.Config, error) {
	certificate, privateKey, creationError := CreateCertificate(template)
	tlsCertificate := new(tls.Certificate)
	tlsCertificate.Certificate = [][]byte{certificate}
	tlsCertificate.PrivateKey = privateKey
	if creationError == nil {

		configuration := new(tls.Config)
		configuration.Certificates = []tls.Certificate{*tlsCertificate}

		return configuration, nil
	}
	log.Print(creationError)
	return nil, creationError
}

func CreateMasterTLSConfiguration() (*tls.Config, error) {
	return CreateTLSConfiguration(CreateMasterSlaveCertificateTemplate())
}

func CreateSlaveTLSConfiguration() (*tls.Config, error) {
	configuration, configurationError := CreateTLSConfiguration(CreateMasterSlaveCertificateTemplate())
	if configurationError == nil {
		configuration.InsecureSkipVerify = true
	}
	return configuration, configurationError
}

func UpgradeServerToTLS(connection net.Conn, configuration *tls.Config) net.Conn {
	return tls.Server(connection, configuration)
}

func UpgradeClientToTLS(connection net.Conn, configuration *tls.Config) net.Conn {
	return tls.Client(connection, configuration)
}

func TLSConnect(address string, configuration *tls.Config) (net.Conn, error) {
	connection, connectionError := net.Dial("tcp", address)
	if connection != nil {
		return UpgradeClientToTLS(connection, configuration), nil
	}
	return nil, connectionError
}
