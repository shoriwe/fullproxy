package socks5

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

// Address represents address in socks protocol.
type Address struct {
	Addr net.IP
	ATYPE
	Port uint16
}

var bufPool = sync.Pool{New: func() interface{} {
	buf := bytes.Buffer{}
	return &buf
}}

// Address return address
// Examples:
//    127.0.0.1:80
//    example.com:443
//    [fe80::1%lo0]:80
func (a *Address) String() string {
	if a.ATYPE == DOMAINNAME {
		return net.JoinHostPort(string(a.Addr), strconv.Itoa(int(a.Port)))
	}
	return net.JoinHostPort(a.Addr.String(), strconv.Itoa(int(a.Port)))
}

var errDomainMaxLengthLimit = errors.New("domain name out of max length")

// Bytes return bytes slice of Address by ver param.
// If ver is socks4, the returned socks4 address format as follows:
//    +----+----+----+----+----+----+....+----+....+----+
//    | DSTPORT |      DSTIP        | USERID       |NULL|
//    +----+----+----+----+----+----+----+----+....+----+
// If ver is socks4 and address type is domain name,
// the returned socks4 address format as follows:
//    +----+----+----+----+----+----+....+----+....+----+....+----+....+----+
//    | DSTPORT |      DSTIP        | USERID       |NULL|   HOSTNAME   |NULL|
//    +----+----+----+----+----+----+----+----+....+----+----+----+....+----+
// If ver is socks5
// the returned socks5 address format as follows:
//    +------+----------+----------+
//    | ATYP | DST.ADDR | DST.PORT |
//    +------+----------+----------+
//    |  1   | Variable |    2     |
//    +------+----------+----------+
// Socks4 call this method return bytes end with NULL, socks4 client use normally,
// Socks4 server should trim terminative NULL.
// Socks4 server should not call this method if server address type is DOMAINNAME
func (a *Address) Bytes(ver VER) ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer buf.Reset()
	defer bufPool.Put(buf)

	// port
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, a.Port)

	switch ver {
	case Version4:
		// socks4a
		buf.Write(port)
		if a.ATYPE == DOMAINNAME {
			buf.Write(net.IPv4(0, 0, 0, 1).To4())
			// NULL
			buf.WriteByte(NULL)
			// hostname
			buf.Write(a.Addr)
		} else if a.ATYPE == IPV4_ADDRESS {
			buf.Write(a.Addr)
		} else {
			return nil, fmt.Errorf("socks4 unsupported address type: %#x", a.ATYPE)
		}
		buf.WriteByte(NULL)
	case Version5:
		// address type
		buf.WriteByte(a.ATYPE)
		// domain name address type
		if a.ATYPE == DOMAINNAME {
			if len(a.Addr) > 255 {
				return nil, errDomainMaxLengthLimit
			}
			buf.WriteByte(byte(len(a.Addr)))
			buf.Write(a.Addr)
		} else if a.ATYPE == IPV4_ADDRESS {
			buf.Write(a.Addr.To4())
		} else if a.ATYPE == IPV6_ADDRESS {
			buf.Write(a.Addr.To16())
		}
		buf.Write(port)
	}

	return buf.Bytes(), nil
}

// readAddress read address info from follows:
//    socks5 server's reply.
//    socks5 client's request.
//    socks5 server's udp reply header.
//    socks5 client's udp request header.
//
//    socks4 client's  request.
//    socks4a client's  request
// exclude: socks4a server's reply, socks4 server's reply. Please use readSocks4ReplyAddress.
func readAddress(r io.Reader, ver VER) (*Address, REP, error) {
	addr := &Address{}

	switch ver {
	case Version4:
		// DST.PORT
		port, err := ReadNBytes(r, 2)
		if err != nil {
			return nil, Rejected, &OpError{Version5, "read", nil, "client dest port", err}
		}
		addr.Port = binary.BigEndian.Uint16(port)
		// DST.IP
		ip, err := ReadNBytes(r, 4)
		if err != nil {
			return nil, Rejected, &OpError{Version4, "read", nil, "\"process request dest ip\"", err}
		}
		addr.ATYPE = IPV4_ADDRESS

		//Discard later bytes until read EOF
		//Please see socks4 request format at(http://ftp.icm.edu.pl/packages/socks/socks4/SOCKS4.protocol)
		_, err = ReadUntilNULL(r)
		if err != nil {
			return nil, Rejected, &OpError{Version4, "read", nil, "\"process request useless header \"", err}
		}

		//Socks4a extension
		//    +----+----+----+----+----+----+----+----+----+----++----++-----+-----++----+
		//    | VN | CD | DSTPORT |      DSTIP        | USERID   |NULL|  HOSTNAME   |NULL|
		//    +----+----+----+----+----+----+----+----+----+----++----++-----+-----++----+
		//       1    1      2              4           variable    1    variable    1
		//The client sets the first three bytes of DSTIP to NULL and
		//the last byte to non-zero. The corresponding IP address is
		//0.0.0.x, where x is non-zero
		if ip[0] == 0 && ip[1] == 0 && ip[2] == 0 &&
			ip[3] != 0 {
			ip, err = ReadUntilNULL(r)
			if err != nil {
				return nil, Rejected, &OpError{Version4, "read", nil, "\"process socks4a extension request\"", err}
			}
			addr.ATYPE = DOMAINNAME
		}
		addr.Addr = ip
		return addr, Granted, nil
	case Version5:
		// ATYP
		aType, err := ReadNBytes(r, 1)
		if err != nil {
			return nil, GENERAL_SOCKS_SERVER_FAILURE, &OpError{Version5, "read", nil, "dest address type", err}
		}
		addr.ATYPE = aType[0]

		var addrLen int
		switch addr.ATYPE {
		case IPV4_ADDRESS:
			addrLen = 4
		case IPV6_ADDRESS:
			addrLen = 16
		case DOMAINNAME:
			fqdnLength, err := ReadNBytes(r, 1)
			if err != nil {
				return nil, GENERAL_SOCKS_SERVER_FAILURE, &OpError{Version5, "read", nil, "\"dest domain name length\"", err}
			}
			addrLen = int(fqdnLength[0])
		default:
			return nil, ADDRESS_TYPE_NOT_SUPPORTED, &OpError{Version5, "", nil, "\"dest address\"", &AtypeError{aType[0]}}
		}

		// DST.ADDR
		ip, err := ReadNBytes(r, addrLen)
		if err != nil {
			return nil, GENERAL_SOCKS_SERVER_FAILURE, err
		}
		addr.Addr = ip

		// DST.PORT
		port, err := ReadNBytes(r, 2)
		if err != nil {
			return nil, GENERAL_SOCKS_SERVER_FAILURE, &OpError{Version5, "read", nil, "client dest port", err}
		}
		addr.Port = binary.BigEndian.Uint16(port)
		return addr, SUCCESSED, nil
	default:
		return nil, UNASSIGNED, &VersionError{ver}
	}
}

// readSocks4ReplyAddress read socks4 reply address. Why don't use readAddress,
// because socks4 reply not end with NULL, they're not compatible
func readSocks4ReplyAddress(r io.Reader, ver VER) (*Address, REP, error) {
	addr := &Address{}

	// DST.PORT
	port, err := ReadNBytes(r, 2)
	if err != nil {
		return nil, Rejected, &OpError{Version5, "read", nil, "client dest port", err}
	}
	addr.Port = binary.BigEndian.Uint16(port)
	// DST.IP
	ip, err := ReadNBytes(r, 4)
	if err != nil {
		return nil, Rejected, &OpError{Version4, "read", nil, "\"process request dest ip\"", err}
	}
	addr.Addr = ip
	addr.ATYPE = IPV4_ADDRESS

	return addr, Granted, nil
}

// UDPAddr return UDP Address.
func (a *Address) UDPAddr() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", a.String())
}

// TCPAddr return TCP Address.
func (a *Address) TCPAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", a.String())
}

// ParseAddress parse address to *Address
// Input Examples:
//    127.0.0.1:80
//    example.com:443
//    [fe80::1%lo0]:80
func ParseAddress(addr string) (*Address, error) {
	Address := new(Address)

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		Address.ATYPE = DOMAINNAME
		Address.Addr = []byte(host)
	} else if ip.To4() != nil {
		Address.ATYPE = IPV4_ADDRESS
		Address.Addr = ip.To4()
	} else if ip.To16() != nil {
		Address.ATYPE = IPV6_ADDRESS
		Address.Addr = ip.To16()
	}
	atoi, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	Address.Port = uint16(atoi)
	return Address, nil
}
