package socks5

import (
	"bytes"
	"errors"
)

// Request The SOCKS request is formed as follows:
//    +----+-----+-------+------+----------+----------+
//    |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
//    +----+-----+-------+------+----------+----------+
//    | 1  |  1  | X'00' |  1   | Variable |    2     |
//    +----+-----+-------+------+----------+----------+
type Request struct {
	VER
	CMD
	RSV uint8
	*Address
}

// UDPHeader Each UDP datagram carries a UDP request
// header with it:
//    +----+------+------+----------+----------+----------+
//    |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
//    +----+------+------+----------+----------+----------+
//    | 2  |  1   |  1   | Variable |    2     | Variable |
//    +----+------+------+----------+----------+----------+
type UDPHeader struct {
	RSV  uint16
	FRAG uint8
	*Address
	Data []byte
}

var errEmptyPayload = errors.New("empty payload")

// PackUDPData add UDP request header before payload.
func PackUDPData(addr *Address, payload []byte) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errEmptyPayload
	}
	if addr == nil {
		return nil, errors.New("addr is nil")
	}
	// RSV, FRAG
	data := []byte{0x00, 0x00, 0x00}
	dest, err := addr.Bytes(Version5)
	if err != nil {
		return nil, err
	}
	// ATYP, DEST.IP, DEST.PORT
	data = append(data, dest...)
	// DATA
	data = append(data, payload...)
	return data, nil
}

// UnpackUDPData split UDP header and payload.
func UnpackUDPData(data []byte) (addr *Address, payload []byte, err error) {
	// trim RSV, FRAG
	data = data[3:]
	buf := bytes.NewBuffer(data)
	addr, _, err = readAddress(buf, Version5)
	if err != nil {
		return nil, nil, err
	}

	payload = buf.Bytes()
	return
}
