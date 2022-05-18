package common

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
)

const (
	SniffSeparator = "\n\n------------------\n\n"
)

func closer(conn1, conn2 io.Closer) {
	_ = conn1.Close()
	_ = conn2.Close()
}

func ForwardTraffic(
	clientConnection net.Conn, targetConnection net.Conn,
	incomingSniffer, outgoingSniffer io.Writer) error {
	defer closer(clientConnection, targetConnection)
	go io.Copy(clientConnection, &ReaderSniffer{
		Writer: incomingSniffer,
		Reader: targetConnection,
	})
	_, err := io.Copy(targetConnection, &ReaderSniffer{
		Writer: outgoingSniffer,
		Reader: clientConnection,
	})
	return err
}

type ReaderSniffer struct {
	Writer io.Writer
	Reader io.Reader
}

func (r *ReaderSniffer) Read(p []byte) (n int, err error) {
	if r.Writer == nil {
		return r.Reader.Read(p)
	}
	length, readReadError := r.Reader.Read(p)
	_, _ = r.Writer.Write(p[:length])
	_, _ = fmt.Fprintf(r.Writer, SniffSeparator)
	return length, readReadError
}

type RequestSniffer struct {
	HeaderDone bool
	Writer     io.Writer
	Request    *http.Request
}

func (r *RequestSniffer) Close() error {
	return r.Request.Body.Close()
}

func (r *RequestSniffer) Read(p []byte) (n int, err error) {
	if r.Writer == nil {
		return r.Request.Body.Read(p)
	}
	if !r.HeaderDone {
		r.HeaderDone = true
		dump, _ := httputil.DumpRequest(r.Request, false)
		_, _ = r.Writer.Write(dump)
	}
	length, readError := r.Request.Body.Read(p)
	_, _ = r.Writer.Write(p[:length])
	_, _ = fmt.Fprintf(r.Writer, SniffSeparator)
	return length, readError
}

type ResponseSniffer struct {
	HeaderDone bool
	Writer     io.Writer
	Response   *http.Response
}

func (r *ResponseSniffer) Close() error {
	return r.Response.Body.Close()
}

func (r *ResponseSniffer) Read(p []byte) (n int, err error) {
	if r.Writer == nil {
		return r.Response.Body.Read(p)
	}
	if !r.HeaderDone {
		r.HeaderDone = true
		dump, _ := httputil.DumpResponse(r.Response, false)
		_, _ = r.Writer.Write(dump)
	}
	length, readError := r.Response.Body.Read(p)
	_, _ = r.Writer.Write(p[:length])
	_, _ = fmt.Fprintf(r.Writer, SniffSeparator)
	return length, readError
}
