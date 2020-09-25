package SOCKS5

const (
	BasicNegotiation byte = 1
	Version          byte = 5
	NoAuthRequired   byte = 0
	InvalidMethod    byte = 1
	UsernamePassword byte = 2
	// SOCKS commands
	Connect      byte = 1
	Bind         byte = 2
	UDPAssociate byte = 3
	// SOCKS valid address types
	IPv4       byte = 1
	DomainName byte = 3
	IPv6       byte = 4
	// SOCKS5 responses
	ConnectionRefused byte = 5
	Succeeded         byte = 0
)

var (
	// SOCKS requests connection responses
	UsernamePasswordSupported = []byte{Version, UsernamePassword}
	NoAuthRequiredSupported   = []byte{Version, NoAuthRequired}
	NoSupportedMethods        = []byte{Version, 255}
	UsernamePasswordSucceededResponse = []byte{1, Succeeded}
	AuthenticationFailed              = []byte{1, 1}
)
