package pipes

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
	"net"
	"time"
)

type Slave struct {
	MasterConnection net.Conn
	NetworkType      string
	MasterC2Address  string
	LoggingMethod    global.LoggingMethod
	IgnoreTrust      bool
	tlsConfig        *tls.Config
}

func (slave *Slave) SetInboundFilter(_ global.IOFilter) error {
	return nil
}

func NewSlave(networkType string, masterC2Address string, loggingMethod global.LoggingMethod, ignoreTrust bool) *Slave {
	return &Slave{
		NetworkType:     networkType,
		MasterC2Address: masterC2Address,
		LoggingMethod:   loggingMethod,
		IgnoreTrust:     ignoreTrust,
	}
}

func (slave *Slave) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	slave.LoggingMethod = loggingMethod
	return nil
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
	global.LogData(slave.LoggingMethod, "Connecting to: ", address)
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
	return ForwardTraffic(clientConnection, targetConnection)
}

func (slave *Slave) command(command byte, clientConnection net.Conn) error {
	switch command {
	case DialCommand:
		return slave.dial(clientConnection)
	}
	_, _ = clientConnection.Write([]byte{UnknownCommand})
	return errors.New("unknown command")
}

func (slave *Slave) serve() error {
	global.LogData(slave.LoggingMethod, "Received client connection from master")
	clientConnection, connectionError := tls.Dial(slave.NetworkType, slave.MasterC2Address, slave.tlsConfig)
	if connectionError != nil {
		return connectionError
	}
	command := make([]byte, 1)
	var bytesReceived int
	bytesReceived, connectionError = clientConnection.Read(command)
	if connectionError != nil {
		return connectionError
	} else if bytesReceived != 1 {
		return errors.New("setup connection error")
	}
	return slave.command(command[0], clientConnection)
}

func (slave *Slave) Serve() error {
	slave.tlsConfig = &tls.Config{
		InsecureSkipVerify: slave.IgnoreTrust,
	}
	global.LogData(slave.LoggingMethod, "Connecting to master at: "+slave.MasterC2Address)
	masterConnection, connectionError := tls.Dial(slave.NetworkType, slave.MasterC2Address, slave.tlsConfig)
	if connectionError != nil {
		return connectionError
	}
	global.LogData(slave.LoggingMethod, "Successfully connected to target at: "+slave.MasterC2Address)
	slave.MasterConnection = masterConnection
	var bytesReceived int
	for {
		bytesReceived, connectionError = slave.MasterConnection.Read(make([]byte, 10))
		if connectionError != nil {
			return connectionError
		}
		for index := 0; index < bytesReceived; index++ {
			go slave.serve()
		}
	}
}
