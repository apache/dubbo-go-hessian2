package exception

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUndeclaredThrowableException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_UndeclaredThrowableException", func(r interface{}) {
		assert.Equal(t, "UndeclaredThrowableException", r.(error).Error())
	})
}
