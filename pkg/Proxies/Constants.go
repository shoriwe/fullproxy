package Proxies

import "github.com/shoriwe/FullProxy/pkg/ConnectionControllers"

func Authenticate(authenticationMethod ConnectionControllers.AuthenticationMethod, username []byte, password []byte) bool {
	if authenticationMethod != nil {
		return authenticationMethod(username, password)
	}
	return true
}
