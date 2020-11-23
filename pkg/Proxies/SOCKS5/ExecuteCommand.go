package SOCKS5

import (
	"bufio"
	"errors"
	"net"
)

func (socks5 *Socks5) ExecuteCommand(
	clientConnection net.Conn, clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetRequestedCommand *byte, targetHostType *byte,
	targetHost *string, targetPort *string) error {

	switch *targetRequestedCommand {
	case Connect:
		return socks5.PrepareConnect(clientConnection, clientConnectionReader, clientConnectionWriter, targetHost, targetPort, targetHostType)
	case Bind:
		return socks5.PrepareBind(clientConnection, clientConnectionReader, clientConnectionWriter, targetHost, targetPort, targetHostType)
	case UDPAssociate:
		return socks5.PrepareUDPAssociate(clientConnection, clientConnectionReader, clientConnectionWriter, targetHost, targetPort, targetHostType)
	}
	return errors.New("unknown command")
}
