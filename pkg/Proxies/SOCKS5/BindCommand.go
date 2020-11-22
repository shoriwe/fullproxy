package SOCKS5

import "errors"

func (socks5 *Socks5)PrepareBind(targetAddress *string,
	targetPort *string, targetAddressType *byte) error {
	_ = socks5.ClientConnection.Close()
	return errors.New("method not implemented yet")
}
