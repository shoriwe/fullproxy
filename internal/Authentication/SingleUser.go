package Authentication

import (
	"bytes"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
)

func SingleUser(username []byte, password []byte) Types.AuthenticationMethod {
	return func(clientUsername []byte, clientPassword []byte) bool {
		return bytes.Equal(username, clientUsername) && bytes.Equal(password, clientPassword)
	}
}
