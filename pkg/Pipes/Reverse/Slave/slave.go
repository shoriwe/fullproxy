package Slave

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"io"
	"net"
)

type Slave struct {
	MasterConnection   net.Conn
	NetworkType        string
	MasterC2Address    string
	MasterProxyAddress string
	LoggingMethod      Types.LoggingMethod
}

func (slave *Slave) SetInboundFilter(_ Types.IOFilter) error {
	panic("inbound rules not supported")
}

func NewSlave(networkType string, masterC2Address string, masterProxyAddress string, loggingMethod Types.LoggingMethod) *Slave {
	return &Slave{NetworkType: networkType, MasterC2Address: masterC2Address, MasterProxyAddress: masterProxyAddress, LoggingMethod: loggingMethod}
}

func (slave *Slave) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	slave.LoggingMethod = loggingMethod
	return nil
}

func (slave *Slave) dial(clientConnection net.Conn) error {
	networkTypeLength := make([]byte, 1)
	bytesReceived, connectionError := clientConnection.Read(networkTypeLength)
	if connectionError != nil {
		return connectionError
	} else if bytesReceived != 1 {
		return errors.New("expecting 1 byte from master")
	}
	rawNetworkType := make([]byte, int(networkTypeLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawNetworkType)
	if connectionError != nil {
		return connectionError
	} else if bytesReceived != int(networkTypeLength[0]) {
		return errors.New("expecting more bytes from master")
	}
	addressLength := make([]byte, 1)
	bytesReceived, connectionError = clientConnection.Read(addressLength)
	if connectionError != nil {
		return connectionError
	} else if bytesReceived != 1 {
		return errors.New("expecting 1 byte from master")
	}
	rawAddress := make([]byte, int(addressLength[0]))
	bytesReceived, connectionError = clientConnection.Read(rawAddress)
	if connectionError != nil {
		return connectionError
	} else if bytesReceived != int(addressLength[0]) {
		return errors.New("expecting more bytes from master")
	}
	networkType := string(rawNetworkType)
	address := string(rawAddress)
	Tools.LogData(slave.LoggingMethod, "Connecting to: ", address)
	var targetConnection net.Conn
	targetConnection, connectionError = net.Dial(networkType, address)
	if connectionError != nil {
		return connectionError
	}
	_, connectionError = clientConnection.Write([]byte{Reverse.NewConnectionSucceeded})
	if connectionError != nil {
		return connectionError
	}
	go io.Copy(targetConnection, clientConnection)
	_, connectionError = io.Copy(clientConnection, targetConnection)
	return connectionError
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

func (slave *Slave) serve(request []byte) error {
	Tools.LogData(slave.LoggingMethod, "Received client connection from master")
	switch request[0] {
	case Reverse.RequestNewMasterConnection:
		break
	default:
		return nil
	}
	clientConnection, connectionError := net.Dial(slave.NetworkType, slave.MasterC2Address)
	if connectionError != nil {
		return connectionError
	}
	defer clientConnection.Close()
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
		request := make([]byte, 1)
		bytesReceived, connectionError = slave.MasterConnection.Read(request)
		if connectionError != nil {
			return connectionError
		} else if bytesReceived != 1 {
			continue
		}
		go slave.serve(request)
	}
}
