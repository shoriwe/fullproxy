package common

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	SniffSeparator    = "\n\n----------\n\n"
	DefaultBufferSize = 1024 * 32
)

func closer(conn1, conn2 io.Closer) {
	_ = conn1.Close()
	_ = conn2.Close()
}

func netCopy(dst, src net.Conn, sniffer io.Writer) error {
	if sniffer == nil {
		_, err := io.Copy(dst, src)
		return err
	}
	var buffer [DefaultBufferSize]byte
	for {
		length, readError := src.Read(buffer[:])
		if readError != nil {
			return readError
		}
		_, writeError := dst.Write(buffer[:length])
		if writeError != nil {
			return writeError
		}
		go sniffer.Write(buffer[:length])
	}

}

func ForwardTraffic(
	clientConnection net.Conn, targetConnection net.Conn,
	incomingSniffer, outgoingSniffer io.Writer) error {
	defer closer(clientConnection, targetConnection)
	go netCopy(clientConnection, targetConnection, incomingSniffer)
	return netCopy(targetConnection, clientConnection, outgoingSniffer)
}

func SniffRequest(w io.Writer, req *http.Request) error {
	if w == nil {
		return nil
	}
	_, _ = fmt.Fprintf(w, SniffSeparator)
	return req.Write(w)
}

func SniffResponse(w io.Writer, res *http.Response) error {
	if w == nil {
		return nil
	}
	_, _ = fmt.Fprintf(w, SniffSeparator)
	return res.Write(w)
}
