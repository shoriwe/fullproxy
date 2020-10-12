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

func CreateCertificateTemplate() *x509.Certificate{
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

func CreateCertificate() ([]byte, *rsa.PrivateKey, error) {
	privateKey, keyGenerationError := rsa.GenerateKey(rand.Reader, 2048)
	if keyGenerationError == nil {
		template := CreateCertificateTemplate()
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

func CreateTLSConfiguration() (*tls.Config, error) {
	certificate, privateKey, creationError := CreateCertificate()
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


func CreateServerTLSConfiguration() (*tls.Config, error) {
	return CreateTLSConfiguration()
}

func CreateClientTLSConfiguration() (*tls.Config, error) {
	configuration, configurationError := CreateTLSConfiguration()
	if configurationError == nil {
		configuration.InsecureSkipVerify = true
	}
	return configuration, configurationError
}

func UpgradeServerToTLS(connection net.Conn, configuration *tls.Config) net.Conn{
	return tls.Server(connection, configuration)
}

func UpgradeClientToTLS(connection net.Conn, configuration *tls.Config) net.Conn {
	return tls.Client(connection, configuration)
}

func TLSConnect(address *string, port *string, configuration *tls.Config) net.Conn {
	connection := Connect(address, port)
	if connection != nil {
		return UpgradeClientToTLS(connection, configuration)
	}
	return nil
}