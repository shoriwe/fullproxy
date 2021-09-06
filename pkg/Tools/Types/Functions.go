package Types

type AuthenticationMethod func(username []byte, password []byte) (bool, error)

type LoggingMethod func(args ...interface{})

type IOFilter func(host string) bool
