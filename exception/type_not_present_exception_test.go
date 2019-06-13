package exception

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeNotPresentException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_TypeNotPresentException", func(r interface{}) {
		assert.Equal(t, "Type exceptiontype1 not present", r.(error).Error())
	})
}