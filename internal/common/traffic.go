package common

import (
	"bytes"
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
	if readReadError != nil {
		return length, readReadError
	}
	_, writeError := r.Writer.Write(p[:length])
	return length, writeError
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
		requestClone := r.Request.Clone(r.Request.Context())
		requestClone.Body = io.NopCloser(bytes.NewReader(nil))
		_ = requestClone.Write(r.Writer)
	}
	length, readError := r.Request.Body.Read(p)
	_, _ = r.Writer.Write(p[:length])
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
		responseClone := &http.Response{
			Status:           r.Response.Status,
			StatusCode:       r.Response.StatusCode,
			Proto:            r.Response.Proto,
			ProtoMajor:       r.Response.ProtoMajor,
			ProtoMinor:       r.Response.ProtoMinor,
			Header:           r.Response.Header,
			Body:             io.NopCloser(bytes.NewReader(nil)),
			ContentLength:    r.Response.ContentLength,
			TransferEncoding: r.Response.TransferEncoding,
			Close:            r.Response.Close,
			Uncompressed:     r.Response.Uncompressed,
			Trailer:          r.Response.Trailer,
			Request:          r.Response.Request,
			TLS:              r.Response.TLS,
		}
		_ = responseClone.Write(r.Writer)
	}
	length, readError := r.Response.Body.Read(p)
	_, _ = r.Writer.Write(p[:length])
	return length, readError
}
