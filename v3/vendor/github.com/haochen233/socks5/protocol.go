package socks5

import "fmt"

type VersionError struct {
	VER
}

func (v *VersionError) Error() string {
	return fmt.Sprintf("error socks protocol version: %d", v.VER)
}

// VER indicates protocol version
type VER = uint8

const (
	Version4 = 0x04
	Version5 = 0x05
)

type MethodError struct {
	METHOD
}

func (m *MethodError) Error() string {
	if _, ok := method2Str[m.METHOD]; ok {
		return fmt.Sprintf("don't support this method %s", method2Str[m.METHOD])
	} else {
		return fmt.Sprintf("unknown mehotd %#x", m.METHOD)
	}
}

// METHOD Defined authentication methods
type METHOD = uint8

const (
	NO_AUTHENTICATION_REQUIRED METHOD = 0x00
	GSSAPI                     METHOD = 0x01
	USERNAME_PASSWORD          METHOD = 0x02
	IANA_ASSIGNED              METHOD = 0x03
	NO_ACCEPTABLE_METHODS      METHOD = 0xff
)

var method2Str = map[METHOD]string{
	NO_AUTHENTICATION_REQUIRED: "NO_AUTHENTICATION_REQUIRED",
	GSSAPI:                     "GSSAPI",
	USERNAME_PASSWORD:          "USERNAME_PASSWORD",
	IANA_ASSIGNED:              "IANA_ASSIGNED",
	NO_ACCEPTABLE_METHODS:      "NO_ACCEPTABLE_METHODS",
}

// CMDError cmd error type
type CMDError struct {
	CMD
}

func (c *CMDError) Error() string {
	if _, ok := cmd2Str[c.CMD]; !ok {
		return fmt.Sprintf("unknown command:%#x", c.CMD)
	}
	return fmt.Sprintf("don't support this command:%s", cmd2Str[c.CMD])
}

// CMD is one of a field in Socks5 Request
type CMD = uint8

const (
	CONNECT       CMD = 0x01
	BIND          CMD = 0x02
	UDP_ASSOCIATE CMD = 0x03
)

var cmd2Str = map[CMD]string{
	CONNECT:       "CONNECT",
	BIND:          "BIND",
	UDP_ASSOCIATE: "UDP_ASSOCIATE",
	Rejected:      "Rejected",
	Granted:       "Granted",
}

type REPError struct {
	REP
}

func (r *REPError) Error() string {
	if _, ok := cmd2Str[r.REP]; !ok {
		return fmt.Sprintf("unknown rep:%#x", r.REP)
	}
	return fmt.Sprintf("don't support this rep:%s", rep2Str[r.REP])
}

// REP is one of a filed in Socks5 Reply
type REP = uint8

//socks5 reply
const (
	SUCCESSED                       REP = 0x00
	GENERAL_SOCKS_SERVER_FAILURE    REP = 0x01
	CONNECTION_NOT_ALLOW_BY_RULESET REP = 0x02
	NETWORK_UNREACHABLE             REP = 0x03
	HOST_UNREACHABLE                REP = 0x04
	CONNECTION_REFUSED              REP = 0x05
	TTL_EXPIRED                     REP = 0x06
	COMMAND_NOT_SUPPORTED           REP = 0x07
	ADDRESS_TYPE_NOT_SUPPORTED      REP = 0x08
	UNASSIGNED                      REP = 0x09
)

var rep2Str = map[REP]string{
	SUCCESSED:                       "successes",
	GENERAL_SOCKS_SERVER_FAILURE:    "general_socks_server_failure",
	CONNECTION_NOT_ALLOW_BY_RULESET: "connection_not_allow_by_ruleset",
	NETWORK_UNREACHABLE:             "network_unreachable",
	HOST_UNREACHABLE:                "host_unreachable",
	CONNECTION_REFUSED:              "connection_refused",
	TTL_EXPIRED:                     "ttl_expired",
	COMMAND_NOT_SUPPORTED:           "command_not_supported",
	ADDRESS_TYPE_NOT_SUPPORTED:      "address_type_not_supported",
	UNASSIGNED:                      "unassigned",
	Granted:                         "granted",
	Rejected:                        "rejected",
}

//socks4 reply
const (
	// Granted means server allow  client request
	Granted = 90
	// Rejected means server refuse client request
	Rejected = 91
)

type AtypeError struct {
	ATYPE
}

func (a *AtypeError) Error() string {
	return fmt.Sprintf("unknown address type:%#x", a.ATYPE)
}

// ATYPE indicates address type in Request and Reply struct
type ATYPE = uint8

const (
	IPV4_ADDRESS ATYPE = 0x01
	DOMAINNAME   ATYPE = 0x03
	IPV6_ADDRESS ATYPE = 0x04
)

var atype2Str = map[ATYPE]string{
	IPV4_ADDRESS: "IPV4_ADDRESS",
	DOMAINNAME:   "DOMAINNAME",
	IPV6_ADDRESS: "IPV6_ADDRESS",
}

const NULL byte = 0
