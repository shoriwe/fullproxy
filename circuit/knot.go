package circuit

import (
	"github.com/shoriwe/fullproxy/v3/utils/network"
)

type Knot interface {
	Next(dial network.DialFunc) (network.CloseFunc, network.DialFunc, error)
}
