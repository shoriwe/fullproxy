package Proxies

import (
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
)

func Authenticate(authenticationMethod Types.AuthenticationMethod, username []byte, password []byte) bool {
	if authenticationMethod != nil {
		return authenticationMethod(username, password)
	}
	return true
}
