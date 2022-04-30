package socks5

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

// GetRandomPort return a random port by specific network.
// The network must be "tcp", "udp".
func GetRandomPort(network string) (err error, port uint16) {
	network = strings.ToLower(network)
	addr := "0.0.0.0:0"

	switch network {
	case "tcp", "tcp4", "tcp6":
		tcpAddr, err := net.ResolveTCPAddr(network, addr)
		if err != nil {
			return err, 0
		}
		ln, err := net.ListenTCP(network, tcpAddr)
		if err != nil {
			return err, 0
		}

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		p, err := strconv.Atoi(portStr)
		port = uint16(p)

		err = ln.Close()
		if err != nil {
			return err, 0
		}
		return err, port
	case "udp", "udp4", "udp6":
		udpAddr, err := net.ResolveUDPAddr(network, addr)
		if err != nil {
			return err, 0
		}
		ln, err := net.ListenUDP(network, udpAddr)
		if err != nil {
			return err, 0
		}

		_, portStr, err := net.SplitHostPort(ln.LocalAddr().String())
		p, err := strconv.Atoi(portStr)
		port = uint16(p)

		err = ln.Close()
		if err != nil {
			return err, 0
		}
		return err, port
	default:
		return errors.New("unknown network type " + network), 0
	}
}

// IsFreePort indicates the port is available.
// The network must be "tcp", "udp".
func IsFreePort(network string, port uint16) error {
	network = strings.ToLower(network)
	portStr := strconv.Itoa(int(port))
	addr := "0.0.0.0:" + portStr

	switch network {
	case "tcp":
		tcpAddr, err := net.ResolveTCPAddr(network, addr)
		if err != nil {
			return err
		}
		ln, err := net.ListenTCP(network, tcpAddr)
		if err != nil {
			return err
		}

		ln.Close()
		if err != nil {
			return err
		}
		return nil
	case "udp":
		udpAddr, err := net.ResolveUDPAddr(network, addr)
		if err != nil {
			return err
		}
		ln, err := net.ListenUDP(network, udpAddr)
		if err != nil {
			return err
		}

		ln.Close()
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unknown network type " + network)
	}
}
