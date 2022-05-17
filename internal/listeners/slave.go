package listeners

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"log"
	"net"
	"time"
)

type Slave struct {
	C2Connection net.Conn
	Config       *tls.Config
}

func (slave *Slave) Accept() (net.Conn, error) {
	var (
		clientConnection net.Conn
	)
	_, connectionError := slave.C2Connection.Read(make([]byte, 1))
	if connectionError != nil {
		return nil, connectionError
	}
	clientConnection, connectionError = tls.Dial(slave.C2Connection.RemoteAddr().Network(), slave.C2Connection.RemoteAddr().String(), slave.Config)
	if connectionError != nil {
		return nil, connectionError
	}
	command := make([]byte, 1)
	var bytesReceived int
	bytesReceived, connectionError = clientConnection.Read(command)
	if connectionError != nil {
		return nil, connectionError
	} else if bytesReceived != 1 {
		return nil, errors.New("setup connection error")
	}
	switch command[0] {
	case DialCommand:
		log.Println(clientConnection.RemoteAddr().String())
		return clientConnection, slave.dial(clientConnection)
	}
	_, _ = clientConnection.Write([]byte{UnknownCommand})
	return nil, errors.New("unknown command")
}

func (slave *Slave) Close() error {
	return slave.C2Connection.Close()
}

func (slave *Slave) Addr() net.Addr {
	return slave.C2Connection.RemoteAddr()
}

func (slave *Slave) dial(clientConnection net.Conn) error {
	networkTypeLength := make([]byte, 1)
	bytesReceived, connectionError := clientConnection.Read(networkTypeLength)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	} else if bytesReceived != 1 {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return errors.New("expecting 1 byte from master")
	}
	rawNetworkType := make([]byte, int(networkTypeLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawNetworkType)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	} else if bytesReceived != int(networkTypeLength[0]) {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return errors.New("expecting more bytes from master")
	}
	addressLength := make([]byte, 1)
	bytesReceived, connectionError = clientConnection.Read(addressLength)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	} else if bytesReceived != 1 {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return errors.New("expecting 1 byte from master")
	}
	rawAddress := make([]byte, int(addressLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawAddress)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	} else if bytesReceived != int(addressLength[0]) {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return errors.New("expecting more bytes from master")
	}
	networkType := string(rawNetworkType)
	address := string(rawAddress)
	var targetConnection net.Conn
	targetConnection, connectionError = net.DialTimeout(networkType, address, time.Minute)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	}
	_, connectionError = clientConnection.Write([]byte{NewConnectionSucceeded})
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{NewConnectionFailed})
		return connectionError
	}
	return common.ForwardTraffic(clientConnection, targetConnection)
}

func (slave *Slave) Serve() error {
	for {
		slave.Accept()
	}
}

func NewSlave(c2Network string, c2Address string, config *tls.Config) (*Slave, error) {
	if config == nil {
		config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	c2Connection, connectionError := tls.Dial(c2Network, c2Address, config)
	if connectionError != nil {
		return nil, connectionError
	}
	return &Slave{
		C2Connection: c2Connection,
		Config:       config,
	}, nil
}
