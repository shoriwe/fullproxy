package Slave

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
	"time"
)

type Slave struct {
	MasterConnection net.Conn
	NetworkType      string
	MasterC2Address  string
	LoggingMethod    Types.LoggingMethod
}

func (slave *Slave) SetInboundFilter(_ Types.IOFilter) error {
	panic("inbound rules not supported")
}

func NewSlave(networkType string, masterC2Address string, loggingMethod Types.LoggingMethod) *Slave {
	return &Slave{NetworkType: networkType, MasterC2Address: masterC2Address, LoggingMethod: loggingMethod}
}

func (slave *Slave) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	slave.LoggingMethod = loggingMethod
	return nil
}

func (slave *Slave) dial(clientConnection net.Conn) error {
	networkTypeLength := make([]byte, 1)
	bytesReceived, connectionError := clientConnection.Read(networkTypeLength)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	} else if bytesReceived != 1 {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return errors.New("expecting 1 byte from master")
	}
	rawNetworkType := make([]byte, int(networkTypeLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawNetworkType)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	} else if bytesReceived != int(networkTypeLength[0]) {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return errors.New("expecting more bytes from master")
	}
	addressLength := make([]byte, 1)
	bytesReceived, connectionError = clientConnection.Read(addressLength)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	} else if bytesReceived != 1 {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return errors.New("expecting 1 byte from master")
	}
	rawAddress := make([]byte, int(addressLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawAddress)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	} else if bytesReceived != int(addressLength[0]) {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return errors.New("expecting more bytes from master")
	}
	networkType := string(rawNetworkType)
	address := string(rawAddress)
	Tools.LogData(slave.LoggingMethod, "Connecting to: ", address)
	var targetConnection net.Conn
	targetConnection, connectionError = net.DialTimeout(networkType, address, time.Minute)
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	}
	_, connectionError = clientConnection.Write([]byte{Reverse.NewConnectionSucceeded})
	if connectionError != nil {
		_, _ = clientConnection.Write([]byte{Reverse.NewConnectionFailed})
		return connectionError
	}
	return Pipes.ForwardTraffic(clientConnection, targetConnection)
}

func (slave *Slave) bind(clientConnection net.Conn) error {
	return nil
}

func (slave *Slave) command(command byte, clientConnection net.Conn) error {
	switch command {
	case Reverse.Dial:
		return slave.dial(clientConnection)
	case Reverse.Bind:
		return slave.bind(clientConnection)
	}
	_, _ = clientConnection.Write([]byte{Reverse.UnknownCommand})
	return errors.New("unknown command")
}

func (slave *Slave) serve() error {
	Tools.LogData(slave.LoggingMethod, "Received client connection from master")
	clientConnection, connectionError := net.DialTimeout(slave.NetworkType, slave.MasterC2Address, time.Minute)
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
	Tools.LogData(slave.LoggingMethod, "Connecting to master at: "+slave.MasterC2Address)
	masterConnection, connectionError := net.Dial(slave.NetworkType, slave.MasterC2Address)
	if connectionError != nil {
		return connectionError
	}
	Tools.LogData(slave.LoggingMethod, "Successfully connected to target at: "+slave.MasterC2Address)
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
