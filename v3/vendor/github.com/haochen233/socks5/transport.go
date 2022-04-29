package socks5

import (
	"io"
	"net"
	"strings"
	"sync"
)

// Transporter transmit data between client and dest server.
type Transporter interface {
	TransportTCP(client *net.TCPConn, remote *net.TCPConn) <-chan error
	TransportUDP(server *UDPConn, request *Request) error
}

type transport struct {
}

const maxLenOfDatagram = 65507

var transportPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, maxLenOfDatagram)
	},
}

// TransportTCP use io.CopyBuffer transmit data.
func (t *transport) TransportTCP(client *net.TCPConn, remote *net.TCPConn) <-chan error {
	errCh := make(chan error)
	var wg = sync.WaitGroup{}

	f := func(dst *net.TCPConn, src *net.TCPConn) {
		defer wg.Done()
		buf := transportPool.Get().([]byte)
		defer transportPool.Put(buf)
		_, err := io.CopyBuffer(dst, src, buf)
		errCh <- err
	}

	wg.Add(2)
	go func() {
		wg.Wait()
		defer client.Close()
		defer remote.Close()
		close(errCh)
	}()

	go f(remote, client)
	go f(client, remote)

	return errCh
}

// TransportUDP forwarding UDP packet between client and dest.
func (t *transport) TransportUDP(server *UDPConn, request *Request) error {
	// Client udp address, limit access to the association.
	clientAddr, err := request.Address.UDPAddr()
	if err != nil {
		return err
	}

	// Record dest address, limit access to the association.
	forwardAddr := make(map[string]struct{})
	buf := transportPool.Get().([]byte)
	defer transportPool.Put(buf)

	defer server.Close()
	for {
		select {
		default:
			// Receive data from remote.
			n, addr, err := server.ReadFromUDP(buf)
			if err != nil {
				return err
			}

			// Should unpack data when data from client.
			if strings.EqualFold(clientAddr.String(), addr.String()) {
				destAddr, payload, err := UnpackUDPData(buf[:n])
				if err != nil {
					return err
				}

				destUDPAddr, err := destAddr.UDPAddr()
				if err != nil {
					return err
				}
				forwardAddr[destUDPAddr.String()] = struct{}{}

				// send payload to dest address
				_, err = server.WriteToUDP(payload, destUDPAddr)
				if err != nil {
					return err
				}
			}

			// Should pack data when data from dest host
			if _, ok := forwardAddr[addr.String()]; ok {
				address, err := ParseAddress(addr.String())
				if err != nil {
					return err
				}

				// packed Data
				packedData, err := PackUDPData(address, buf[:n])
				if err != nil {
					return err
				}

				// send payload to client
				_, err = server.WriteToUDP(packedData, clientAddr)
				if err != nil {
					return err
				}
			}
		case <-server.CloseCh():
			return nil
		}
	}
}

var DefaultTransporter Transporter = &transport{}
