package exception

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMalformedParameterizedTypeException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_MalformedParameterizedTypeException", func(r interface{}) {
		assert.Equal(t, "MalformedParameterizedType", r.(error).Error())
	})
}