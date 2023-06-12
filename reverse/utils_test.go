package reverse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailResponse(t *testing.T) {
	assert.NotNil(t, FailResponse(nil))
}
