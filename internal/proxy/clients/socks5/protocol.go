package socks5

import (
	"encoding/binary"
	"errors"
	socks52 "github.com/shoriwe/fullproxy/v3/internal/proxy/servers/socks5"
	"net"
	"strconv"
)

type Socks5 struct {
	Network, Address   string
	Username, Password string
}

func prepareAddress(address string) (addressType byte, addressBytes []byte, portBytes []byte, err error) {
	host, port, splitError := net.SplitHostPort(address)
	if splitError != nil {
		return 0, nil, nil, splitError
	}
	if ip := net.ParseIP(host); ip != nil {
		if ipv4 := ip.To4(); ipv4 != nil {
			addressType = socks52.IPv4
			addressBytes = ipv4
		} else if ipv6 := ip.To16(); ipv6 != nil {
			addressType = socks52.IPv6
			addressBytes = ipv6
		} else {
			return 0, nil, nil, socks52.UnknownAddressType
		}
	} else {
		addressType = socks52.DomainName
		addressBytes = append([]byte{byte(len(host))}, []byte(host)...)
	}
	portValue, parseError := strconv.Atoi(port)
	if parseError != nil {
		return 0, nil, nil, parseError
	}
	portBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes[:], uint16(portValue))
	return addressType, addressBytes, portBytes, nil
}

func (socks5 *Socks5) authenticate() (net.Conn, error) {
	proxyConnection, connectionError := net.Dial(socks5.Network, socks5.Address)
	if connectionError != nil {
		return nil, connectionError
	}
	if socks5.Username == "" || socks5.Password == "" {
		_, writeError := proxyConnection.Write([]byte{socks52.SocksV5, 1, socks52.NoAuthRequired})
		if writeError != nil {
			return nil, writeError
		}
		var authResponse [2]byte
		_, readError := proxyConnection.Read(authResponse[:])
		if readError != nil {
			return nil, readError
		}
		if authResponse[0] != socks52.SocksV5 {
			return nil, socks52.SocksVersionNotSupported
		}
		if authResponse[1] != socks52.NoAuthRequired {
			return nil, socks52.UnsupportedAuthenticationMethod
		}
	} else {
		_, writeError := proxyConnection.Write([]byte{socks52.SocksV5, 2, socks52.NoAuthRequired, socks52.UsernamePassword})
		if writeError != nil {
			return nil, writeError
		}
		var authResponse [2]byte
		_, readError := proxyConnection.Read(authResponse[:])
		if readError != nil {
			return nil, readError
		}
		if authResponse[0] != socks52.SocksV5 {
			return nil, socks52.SocksVersionNotSupported
		}
		switch authResponse[1] {
		case socks52.NoAuthRequired:
			break
		case socks52.UsernamePassword:
			_, writeError = proxyConnection.Write(
				append(
					[]byte{
						socks52.BasicNegotiation,
						byte(len(socks5.Username)),
					},
					append(
						[]byte(socks5.Username),
						append(
							[]byte{byte(len(socks5.Password))},
							[]byte(socks5.Password)...,
						)...,
					)...,
				),
			)
			if writeError != nil {
				return nil, writeError
			}
			_, readError = proxyConnection.Read(authResponse[:])
			if readError != nil {
				return nil, readError
			}
			if authResponse[0] != socks52.BasicNegotiation || authResponse[1] != socks52.SucceedAuthentication {
				return nil, socks52.SocksVersionNotSupported
			}
		default:
			return nil, socks52.UnsupportedAuthenticationMethod
		}
	}

	return proxyConnection, nil
}

func (socks5 *Socks5) Dial(_, address string) (net.Conn, error) {
	proxyConnection, connectionError := socks5.authenticate()
	if connectionError != nil {
		return nil, connectionError
	}

	addressType, addressBytes, portBytes, parseError := prepareAddress(address)
	if parseError != nil {
		return nil, parseError
	}
	_, writeError := proxyConnection.Write(
		append(
			append(
				[]byte{socks52.SocksV5, socks52.Connect, 0x00, addressType},
				addressBytes...,
			),
			portBytes[0],
			portBytes[1],
		),
	)
	if writeError != nil {
		return nil, writeError
	}
	var reply [50]byte
	_, readError := proxyConnection.Read(reply[:])
	if readError != nil {
		return nil, readError
	}
	if reply[1] != socks52.ConnectionSucceed {
		return nil, errors.New("connection failed")
	}
	return proxyConnection, nil
}

type Listener struct {
	ProxyConnection net.Conn
	Address         net.Addr
}

type AcceptConnection struct {
	ProxyConnection net.Conn
}

func (l *Listener) Accept() (net.Conn, error) {
	var reply [50]byte
	_, readError := l.ProxyConnection.Read(reply[:])
	if readError != nil {
		return nil, readError
	}
	return l.ProxyConnection, nil
}

func (l *Listener) Close() error {
	return l.ProxyConnection.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.Address
}

func (socks5 *Socks5) Listen(_, address string) (net.Listener, error) {
	proxyConnection, connectionError := socks5.authenticate()
	if connectionError != nil {
		return nil, connectionError
	}

	addressType, addressBytes, portBytes, parseError := prepareAddress(address)
	if parseError != nil {
		return nil, parseError
	}
	_, writeError := proxyConnection.Write(
		append(
			append(
				[]byte{socks52.SocksV5, socks52.Bind, 0x00, addressType},
				addressBytes...,
			),
			portBytes[0],
			portBytes[1],
		),
	)
	if writeError != nil {
		return nil, writeError
	}
	var reply [50]byte
	_, readError := proxyConnection.Read(reply[:])
	if readError != nil {
		return nil, readError
	}
	if reply[1] != socks52.ConnectionSucceed {
		return nil, errors.New("connection failed")
	}
	var listenAddress net.Addr
	listenAddress, parseError = parseReply(reply)
	if parseError != nil {
		return nil, parseError
	}
	return &Listener{ProxyConnection: proxyConnection, Address: listenAddress}, nil
}

func parseReply(reply [50]byte) (net.Addr, error) {
	var (
		dstRawAddress, dstRawPort []byte
		dstAddress                net.IP
		dstPort                   int
	)
	switch reply[3] {
	case socks52.IPv4:
		dstRawAddress = reply[4 : 4+4]
		dstRawPort = reply[4+4 : 4+4+2]

		dstAddress = net.IP(dstRawAddress).To4()
	case socks52.IPv6:
		dstRawAddress = reply[4 : 4+16]
		dstRawPort = reply[4+16 : 4+16+2]

		dstAddress = net.IP(dstRawAddress).To16()
	default:
		return nil, socks52.UnknownAddressType
	}
	dstPort = int(binary.BigEndian.Uint16(dstRawPort))
	return &net.TCPAddr{
		IP:   dstAddress,
		Port: dstPort,
	}, nil
}
