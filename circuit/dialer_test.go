package circuit

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDialer_Dial(t *testing.T) {
	d := &Dialer{DialFunc: net.Dial}
	_, err := d.Dial("tcp", "1111111111111111111111")
	assert.NotNil(t, err)
}
