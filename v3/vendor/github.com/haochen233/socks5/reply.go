package socks5

// Reply a reply formed as follows:
//    +----+-----+-------+------+----------+----------+
//    |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
//    +----+-----+-------+------+----------+----------+
//    | 1  |  1  | X'00' |  1   | Variable |    2     |
//    +----+-----+-------+------+----------+----------+
type Reply struct {
	VER
	REP
	RSV uint8
	*Address
}
