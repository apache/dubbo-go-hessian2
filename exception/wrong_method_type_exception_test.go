package exception
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrongMethodTypeException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_WrongMethodTypeException", func(r interface{}) {
		assert.Equal(t, "WrongMethodTypeException", r.(error).Error())
	})
}
