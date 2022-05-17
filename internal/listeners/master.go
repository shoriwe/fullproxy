package listeners

import (
	"crypto/tls"
	"errors"
	"github.com/shoriwe/fullproxy/v3/internal/common"
	"net"
)

// Master
/*
C2Address is the address where the slave is going to connect to.
ProxyAddress is the address where the clients are going to connect to.
*/
type Master struct {
	filters         Filters
	SlaveConnection *tls.Conn
	C2Listener      net.Listener
	ProxyListener   net.Listener
}

func (master *Master) Filter() Filters {
	return master.filters
}

func (master *Master) SetFilters(filters Filters) {
	master.filters = filters
}

func (master *Master) Accept() (net.Conn, error) {
	return master.ProxyListener.Accept()
}

func (master *Master) Close() error {
	_ = master.SlaveConnection.Close()
	_ = master.C2Listener.Close()
	return master.ProxyListener.Close()
}

func (master *Master) Addr() net.Addr {
	return master.ProxyListener.Addr()
}

func (master *Master) Dial(network, address string) (net.Conn, error) {
	numberOfBytesWritten, connectionError := master.SlaveConnection.Write([]byte{RequestNewMasterConnectionCommand})
	if connectionError != nil {
		return nil, connectionError
	} else if numberOfBytesWritten != 1 {
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

func (master *Master) Init() error {
	slaveConnection, connectionError := master.C2Listener.Accept()
	if connectionError != nil {
		return connectionError
	}
	master.SlaveConnection = slaveConnection.(*tls.Conn)
	return master.SlaveConnection.Handshake()
}

func NewMaster(
	network, address string,
	config *tls.Config,
	c2Network, c2Address string,
	c2Config *tls.Config,
) (Listener, error) {
	var (
		proxyListener, c2Listener net.Listener
		listenError               error
	)
	if config == nil {
		proxyListener, listenError = net.Listen(network, address)
	} else {
		proxyListener, listenError = tls.Listen(network, address, config)
	}
	if listenError != nil {
		return nil, listenError
	}
	if c2Config == nil {
		certificates, genCertsError := common.SelfSignCertificate()
		if genCertsError != nil {
			return nil, genCertsError
		}
		c2Config = &tls.Config{
			Certificates: certificates,
		}
	}
	c2Listener, listenError = tls.Listen(c2Network, c2Address, c2Config)
	if listenError != nil {
		return nil, listenError
	}
	return &Master{
		ProxyListener: proxyListener,
		C2Listener:    c2Listener,
	}, nil
}
