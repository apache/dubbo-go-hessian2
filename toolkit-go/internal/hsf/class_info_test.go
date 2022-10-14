package hsf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveGenericName(t *testing.T) {
	s := removeGenericName("122<fafa,>")
	assert.Equal(t, s, "122")
}
