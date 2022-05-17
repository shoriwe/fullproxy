package listeners

import "net"

type TCPListener struct {
	Listener *net.TCPListener
	Filters  Filters
}

func (T *TCPListener) FilterAccept(address string) error {
	return T.Filters.Accept(address)
}

func (T *TCPListener) Accept() (net.Conn, error) {
	connection, connectionError := T.Listener.Accept()
	if connectionError != nil {
		return nil, connectionError
	}
	filterError := T.FilterAccept(connection.RemoteAddr().String())
	if filterError != nil {
		_ = connection.Close()
		return nil, filterError
	}
	return connection, nil
}

func (T *TCPListener) Close() error {
	return T.Listener.Close()
}

func (T *TCPListener) Addr() net.Addr {
	return T.Listener.Addr()
}
