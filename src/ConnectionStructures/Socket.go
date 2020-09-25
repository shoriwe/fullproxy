package ConnectionStructures

import (
	"bufio"
)

type SocketReader interface {
	Read(buffer []byte) (int, error)
}

type BasicConnectionReader struct {
	Reader *bufio.Reader
}

func (connection BasicConnectionReader) Read(buffer []byte) (int, error) {
	return connection.Reader.Read(buffer)
}

type SocketWriter interface {
	Write(buffer []byte) (int, error)
	Flush() error
}

type BasicConnectionWriter struct {
	Writer *bufio.Writer
}

func (connection BasicConnectionWriter) Write(buffer []byte) (int, error) {
	return connection.Writer.Write(buffer)
}

func (connection BasicConnectionWriter) Flush() error {
	return connection.Writer.Flush()
}
