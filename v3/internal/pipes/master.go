package pipes

import (
	"errors"
	"github.com/shoriwe/FullProxy/v3/internal/global"
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
	LoggingMethod    global.LoggingMethod
	InboundFilter    global.IOFilter
	Protocol         global.ProxyProtocol
}

func NewMaster(networkType string, c2Address string, proxyAddress string, loggingMethod global.LoggingMethod, inboundFilter global.IOFilter, protocol global.ProxyProtocol) *Master {
	return &Master{NetworkType: networkType, C2Address: c2Address, ProxyAddress: proxyAddress, LoggingMethod: loggingMethod, InboundFilter: inboundFilter, Protocol: protocol}
}

func (master *Master) SetInboundFilter(filter global.IOFilter) error {
	master.InboundFilter = filter
	return nil
}

func (master *Master) SetLoggingMethod(loggingMethod global.LoggingMethod) error {
	master.LoggingMethod = loggingMethod
	return nil
}

func (master *Master) protocolDialFunc() global.DialFunc {
	return func(network, address string) (net.Conn, error) {
		numberOfBytesWritten, connectionError := master.MasterConnection.Write([]byte{RequestNewMasterConnectionCommand})
		if connectionError != nil {
			master.finish = true
			return nil, connectionError
		} else if numberOfBytesWritten != 1 {
			master.finish = true
			return nil, errors.New("protocol error")
		}
		var targetConnection net.Conn
		targetConnection, connectionError = master.C2Listener.Accept()
		if connectionError != nil {
			return nil, connectionError
		}
		// Request connection to target
		var request []byte
		networkLength := len(network)
		addressLength := len(address)
		payloadLength := 3 + networkLength + addressLength
		request = append(request, DialCommand)
		request = append(request, byte(networkLength))
		request = append(request, []byte(network)...)
		request = append(request, byte(addressLength))
		request = append(request, []byte(address)...)
		var bytesWritten int
		bytesWritten, connectionError = targetConnection.Write(request)
		if connectionError != nil {
			_ = targetConnection.Close()
			return nil, connectionError
		} else if bytesWritten != payloadLength {
			_ = targetConnection.Close()
			return nil, errors.New("new connection request error")
		}
		response := make([]byte, 1)
		var bytesReceived int
		bytesReceived, connectionError = targetConnection.Read(response)
		if connectionError != nil {
			_ = targetConnection.Close()
			return nil, connectionError
		} else if bytesReceived != 1 {
			_ = targetConnection.Close()
			return nil, errors.New("new connection request error")
		}
		switch response[0] {
		case NewConnectionSucceeded:
			return targetConnection, nil
		}
		return nil, errors.New("new connection request error")
	}
}

func (master *Master) serve(client net.Conn) error {
	global.LogData(master.LoggingMethod, "Received connection from: ", client.RemoteAddr().String())
	if !global.FilterInbound(master.InboundFilter, global.ParseIP(client.RemoteAddr().String()).String()) {
		return errors.New("Connection denied!")
	}
	return master.Protocol.Handle(client)
}

func (master *Master) Serve() error {
	master.finish = false
	global.LogData(master.LoggingMethod, "Listening at: "+master.C2Address)
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
	global.LogData(master.LoggingMethod, "Waiting for slave to connect")
	slaveConnection, connectionError := master.C2Listener.Accept()
	if connectionError != nil {
		return connectionError
	}
	global.LogData(master.LoggingMethod, "slave Address: "+slaveConnection.RemoteAddr().String())
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
