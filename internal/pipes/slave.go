package pipes

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"net"
	"time"
)

type Slave struct {
	MasterConnection              net.Conn
	NetworkType                   string
	MasterC2Address               string
	LoggingMethod                 LoggingMethod
	InboundFilter, OutboundFilter IOFilter
	IgnoreTrust                   bool
	tlsConfig                     *tls.Config
}

func (slave *Slave) FilterInbound(addr net.Addr) error {
	if slave.InboundFilter != nil {
		return slave.InboundFilter(addr)
	}
	return nil
}

func (slave *Slave) FilterOutbound(addr net.Addr) error {
	if slave.OutboundFilter != nil {
		return slave.OutboundFilter(addr)
	}
	return nil
}

func (slave *Slave) LogData(a ...interface{}) {
	if slave.LoggingMethod != nil {
		slave.LoggingMethod(a...)
	}
}

func (slave *Slave) SetOutboundFilter(filter IOFilter) {
	slave.OutboundFilter = filter
}

func (slave *Slave) SetInboundFilter(filter IOFilter) {
	slave.InboundFilter = filter
}

func (slave *Slave) SetLoggingMethod(loggingMethod LoggingMethod) {
	slave.LoggingMethod = loggingMethod
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
	slave.LogData("Connecting to: ", address)
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

func (slave *Slave) command(command byte, clientConnection net.Conn) error {
	switch command {
	case DialCommand:
		return slave.dial(clientConnection)
	}
	_, _ = clientConnection.Write([]byte{UnknownCommand})
	return errors.New("unknown command")
}

func (slave *Slave) serve() error {
	slave.LogData("Received client connection from master")
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
	slave.LogData("Connecting to master at: " + slave.MasterC2Address)
	masterConnection, connectionError := tls.Dial(slave.NetworkType, slave.MasterC2Address, slave.tlsConfig)
	if connectionError != nil {
		return connectionError
	}
	slave.LogData("Successfully connected to target at: " + slave.MasterC2Address)
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

func NewSlave(networkType string, masterC2Address string, loggingMethod LoggingMethod, ignoreTrust bool) Pipe {
	return &Slave{
		NetworkType:     networkType,
		MasterC2Address: masterC2Address,
		LoggingMethod:   loggingMethod,
		IgnoreTrust:     ignoreTrust,
	}
}
