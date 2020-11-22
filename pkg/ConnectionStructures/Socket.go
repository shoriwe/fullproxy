package ConnectionStructures

import (
	"bufio"
	"net"
)

func CreateSocketConnectionReaderWriter(connection net.Conn) (*bufio.Reader, *bufio.Writer) {
	return bufio.NewReader(connection), bufio.NewWriter(connection)
}
