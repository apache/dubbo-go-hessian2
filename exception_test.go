package hessian

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThrowable(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_throwable", func(r interface{}) {
		assert.Equal(t, "exception", r.(error).Error())
	})
}

func TestException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_exception", func(r interface{}) {
		assert.Equal(t, "exception", r.(error).Error())
	})
}
