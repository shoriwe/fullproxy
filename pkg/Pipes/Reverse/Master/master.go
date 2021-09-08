package Master

import (
	"errors"
	"fmt"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
)

type Master struct {
	MasterConnection net.Conn
	C2Listener       net.Listener
	ProxyListener    net.Listener
	LoggingMethod    Types.LoggingMethod
	InboundFilter    Types.IOFilter
	Protocol         Types.ProxyProtocol
}

func NewMaster(host string, c2Port string, proxyPort string, loggingMethod Types.LoggingMethod, inboundFilter Types.IOFilter, protocol Types.ProxyProtocol) (*Master, error) {
	var (
		c2Listener    net.Listener
		proxyListener net.Listener
		listenError   error
	)
	c2Listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", host, c2Port))
	if listenError != nil {
		return nil, listenError
	}
	proxyListener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", host, proxyPort))
	if listenError != nil {
		return nil, listenError
	}
	return &Master{C2Listener: c2Listener, ProxyListener: proxyListener, LoggingMethod: loggingMethod, InboundFilter: inboundFilter, Protocol: protocol}, nil
}

func (master *Master) SetInboundFilter(filter Types.IOFilter) error {
	master.InboundFilter = filter
	return nil
}

func (master *Master) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	master.LoggingMethod = loggingMethod
	return nil
}

func (master *Master) serve(client net.Conn) {
	defer client.Close()
	master.Protocol.Handle(client)
}

func (master *Master) protocolDialFunc() Types.DialFunc {
	return func(network, address string) (net.Conn, error) {
		_, connectionError := master.MasterConnection.Write([]byte{Reverse.RequestNewMasterConnection})
		if connectionError != nil {
			return nil, connectionError
		}
		var c2Connection net.Conn
		c2Connection, connectionError = master.C2Listener.Accept()
		if connectionError != nil {
			return nil, connectionError
		}
		// Request connection to target
		var request []byte
		networkLength := len(network)
		addressLength := len(address)
		payloadLength := 2 + networkLength + 1 + addressLength
		request = append(request, Reverse.Dial)
		request = append(request, byte(networkLength))
		request = append(request, []byte(network)...)
		request = append(request, byte(addressLength))
		request = append(request, []byte(address)...)
		var bytesWritten int
		bytesWritten, connectionError = c2Connection.Write(request)
		if connectionError != nil {
			_ = c2Connection.Close()
			return nil, connectionError
		} else if bytesWritten != payloadLength {
			_ = c2Connection.Close()
			return nil, errors.New("new connection request error")
		}
		response := make([]byte, 1)
		var bytesReceived int
		bytesReceived, connectionError = c2Connection.Read(response)
		if connectionError != nil {
			_ = c2Connection.Close()
			return nil, connectionError
		} else if bytesReceived != 1 {
			_ = c2Connection.Close()
			return nil, errors.New("new connection request error")
		}
		switch response[0] {
		case Reverse.NewConnectionSucceeded:
			return c2Connection, nil
		}
		return nil, errors.New("new connection request error")
	}
}

func (master *Master) Serve() error {
	slaveConnection, connectionError := master.C2Listener.Accept()
	if connectionError != nil {
		return connectionError
	}
	// defer master.ProxyListener.Close()
	// defer master.C2Listener.Close()
	master.MasterConnection = slaveConnection
	master.Protocol.SetDial(master.protocolDialFunc())
	var clientConnection net.Conn
	for {
		clientConnection, connectionError = master.ProxyListener.Accept()
		if connectionError != nil {
			return connectionError
		}
		go master.serve(clientConnection)
	}
}
