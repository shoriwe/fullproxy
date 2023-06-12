package compose

import "github.com/things-go/go-socks5"

type AuthMethod struct {
	Raw map[string]string
}

func (am *AuthMethod) Socks5() (auth socks5.Authenticator, err error) {
	switch {
	case am.Raw != nil:
		auth = socks5.UserPassAuthenticator{
			Credentials: socks5.StaticCredentials(am.Raw),
		}
	default:
		auth = nil
	}
	return auth, err
}

type AuthMethods []AuthMethod

func (am AuthMethods) Socks5() ([]socks5.Authenticator, error) {
	authMethods := make([]socks5.Authenticator, 0, len(am))
	for _, a := range am {
		auth, err := a.Socks5()
		if err != nil {
			return nil, err
		}
		if auth != nil {
			authMethods = append(authMethods, auth)
		}
	}
	return authMethods, nil
}
