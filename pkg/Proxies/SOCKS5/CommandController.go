package SOCKS5

import (
	"errors"
)

func (socks5 *Socks5)HandleCommandExecution(
	targetRequestedCommand *byte, targetAddressType *byte,
	targetAddress *string, targetPort *string) error {

	switch *targetRequestedCommand {
	case Connect:
		return socks5.PrepareConnect(targetAddress, targetPort, targetAddressType)
	case Bind:
		return socks5.PrepareBind(targetAddress, targetPort, targetAddressType)
	case UDPAssociate:
		return socks5.PrepareUDPAssociate(targetAddress, targetPort, targetAddressType)
	}
	return errors.New("unknown command")
}
