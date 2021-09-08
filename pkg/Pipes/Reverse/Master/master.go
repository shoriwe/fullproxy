package Master

import (
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Pipes/Reverse"
	"github.com/shoriwe/FullProxy/pkg/Tools"
	"github.com/shoriwe/FullProxy/pkg/Tools/Types"
	"net"
)

type Master struct {
	finish           bool
	NetworkType      string
	C2Address        string
	ProxyAddress     string
	MasterConnection net.Conn
	C2Listener       net.Listener
	ProxyListener    net.Listener
	LoggingMethod    Types.LoggingMethod
	InboundFilter    Types.IOFilter
	Protocol         Types.ProxyProtocol
}

func NewMaster(networkType string, c2Address string, proxyAddress string, loggingMethod Types.LoggingMethod, inboundFilter Types.IOFilter, protocol Types.ProxyProtocol) *Master {
	return &Master{NetworkType: networkType, C2Address: c2Address, ProxyAddress: proxyAddress, LoggingMethod: loggingMethod, InboundFilter: inboundFilter, Protocol: protocol}
}

func (master *Master) SetInboundFilter(filter Types.IOFilter) error {
	master.InboundFilter = filter
	return nil
}

func (master *Master) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	master.LoggingMethod = loggingMethod
	return nil
}

func (master *Master) serve(client net.Conn) error {
	Tools.LogData(master.LoggingMethod, "Received connection from: ", client.RemoteAddr().String())
	if !Tools.FilterInbound(master.InboundFilter, net.ParseIP(client.RemoteAddr().String()).String()) {
		return errors.New("Connection denied!")
	}
	defer client.Close()
	return master.Protocol.Handle(client)
}

func (master *Master) protocolDialFunc() Types.DialFunc {
	return func(network, address string) (net.Conn, error) {
		_, connectionError := master.MasterConnection.Write([]byte{Reverse.RequestNewMasterConnection})
		if connectionError != nil {
			master.finish = true
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
	master.finish = false
	Tools.LogData(master.LoggingMethod, "Listening at: "+master.C2Address)
	c2Listener, listenError := net.Listen(master.NetworkType, master.C2Address)
	if listenError != nil {
		return listenError
	}
	master.C2Listener = c2Listener
	defer master.C2Listener.Close()
	var proxyListener net.Listener
	proxyListener, listenError = net.Listen(master.NetworkType, master.ProxyAddress)
	if listenError != nil {
		return listenError
	}
	master.ProxyListener = proxyListener
	defer master.ProxyListener.Close()
	//
	Tools.LogData(master.LoggingMethod, "Waiting for slave to connect")
	slaveConnection, connectionError := master.C2Listener.Accept()
	if connectionError != nil {
		return connectionError
	}
	Tools.LogData(master.LoggingMethod, "Slave Address: "+slaveConnection.RemoteAddr().String())
	master.MasterConnection = slaveConnection
	master.Protocol.SetDial(master.protocolDialFunc())
	var clientConnection net.Conn
	for !master.finish {
		clientConnection, connectionError = master.ProxyListener.Accept()
		if connectionError != nil {
			return connectionError
		}
		go master.serve(clientConnection)
	}
	return nil
}
