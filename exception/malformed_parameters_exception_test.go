package exception
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMalformedParametersException(t *testing.T) {
	testDecodeFrameworkFunc(t, "throw_MalformedParametersException", func(r interface{}) {
		assert.Equal(t, "MalformedParametersException", r.(error).Error())
	})
}