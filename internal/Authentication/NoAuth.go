package Authentication

func NoAuthentication(_ []byte, _ []byte) bool {
	return true
}
