package pipes

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"github.com/shoriwe/fullproxy/v3/internal/proxy/servers"
	"net"
)

// Master
/*
C2Address is the address where the slave is going to connect to.
ProxyAddress is the address where the clients are going to connect to.
*/
type Master struct {
	finish                        bool
	NetworkType                   string
	C2Address                     string
	ProxyAddress                  string
	SlaveConnection               net.Conn
	C2Listener                    net.Listener
	ProxyListener                 net.Listener
	LoggingMethod                 LoggingMethod
	InboundFilter, OutboundFilter IOFilter
	Protocol                      servers.Protocol
	TLSConfig                     *tls.Config
}

func (master *Master) Dial(network, address string) (net.Conn, error) {
	if filterError := master.FilterOutbound(address); filterError != nil {
		return nil, filterError
	}
	numberOfBytesWritten, connectionError := master.SlaveConnection.Write([]byte{RequestNewMasterConnectionCommand})
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
	networkLength := len(network)
	addressLength := len(address)
	payloadLength := 3 + networkLength + addressLength
	request := make([]byte, 0, payloadLength)
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
		return nil, SlaveConnectionRequestError
	}
	response := make([]byte, 1)
	var bytesReceived int
	bytesReceived, connectionError = targetConnection.Read(response)
	if connectionError != nil {
		_ = targetConnection.Close()
		return nil, connectionError
	} else if bytesReceived != 1 {
		_ = targetConnection.Close()
		return nil, SlaveConnectionRequestError
	}
	switch response[0] {
	case NewConnectionSucceeded:
		return targetConnection, nil
	}
	return nil, SlaveConnectionRequestError
}

func (master *Master) Listen(_, _ string) (net.Listener, error) {
	return nil, errors.New("not supported for master/slave protocol")
}

func (master *Master) LogData(a ...interface{}) {
	if master.LoggingMethod != nil {
		master.LoggingMethod(a...)
	}
}

func (master *Master) SetOutboundFilter(filter IOFilter) {
	master.OutboundFilter = filter
}

func (master *Master) FilterInbound(address string) error {
	if master.InboundFilter != nil {
		return master.InboundFilter(address)
	}
	return nil
}

func (master *Master) FilterOutbound(address string) error {
	if master.OutboundFilter != nil {
		return master.OutboundFilter(address)
	}
	return nil
}

func (master *Master) SetInboundFilter(filter IOFilter) {
	master.InboundFilter = filter
}

func (master *Master) SetLoggingMethod(loggingMethod LoggingMethod) {
	master.LoggingMethod = loggingMethod
}

func (master *Master) serve(client net.Conn) {
	master.LogData("Received connection from: ", client.RemoteAddr().String())
	if err := master.FilterInbound(client.RemoteAddr().String()); err != nil {
		master.LogData(err)
		return
	}
	handleError := master.Protocol.Handle(client)
	if handleError != nil {
		master.LogData(handleError)
	}
}

func (master *Master) Serve() error {
	if master.TLSConfig.Certificates == nil {
		var selfSignedError error
		master.TLSConfig.Certificates, selfSignedError = common.SelfSignCertificate()
		if selfSignedError != nil {
			return selfSignedError
		}
	}
	master.finish = false
	master.LogData("Listening at: " + master.C2Address)
	c2Listener, listenError := tls.Listen(master.NetworkType, master.C2Address, master.TLSConfig)
	if listenError != nil {
		return listenError
	}
	master.C2Listener = c2Listener
	defer func(C2Listener net.Listener) {
		err := C2Listener.Close()
		if err != nil {
			master.LogData(err)
		}
	}(master.C2Listener)
	var proxyListener net.Listener
	proxyListener, listenError = net.Listen(master.NetworkType, master.ProxyAddress)
	if listenError != nil {
		return listenError
	}
	master.ProxyListener = proxyListener
	defer func(ProxyListener net.Listener) {
		err := ProxyListener.Close()
		if err != nil {
			master.LogData(err)
		}
	}(master.ProxyListener)
	//
	master.LogData("Waiting for slave to connect")
	slaveConnection, connectionError := master.C2Listener.Accept()
	if connectionError != nil {
		return connectionError
	}
	master.LogData("slave Address: " + slaveConnection.RemoteAddr().String())
	master.SlaveConnection = slaveConnection
	master.Protocol.SetDial(master.Dial)
	master.Protocol.SetListen(master.Listen)
	master.Protocol.SetListenAddress(master.ProxyListener.Addr())
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

func NewMaster(
	networkType, c2Address, proxyAddress string,
	loggingMethod LoggingMethod,
	inboundFilter, outboundFilter IOFilter,
	protocol servers.Protocol,
	certificates []tls.Certificate,
) Pipe {
	return &Master{
		NetworkType:    networkType,
		C2Address:      c2Address,
		ProxyAddress:   proxyAddress,
		LoggingMethod:  loggingMethod,
		InboundFilter:  inboundFilter,
		OutboundFilter: outboundFilter,
		Protocol:       protocol,
		TLSConfig:      &tls.Config{Certificates: certificates},
	}
}
