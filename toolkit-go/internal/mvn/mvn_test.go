package mvn

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJarLocalPath(t *testing.T) {
	res, err := os.ReadDir("/Users/seven/.toolkit/.m2/repository/com/autonavi/aos/tmp/amap-aos-tmp-order-kernel-api")
	assert.Nil(t, err)
	for _, re := range res {
		fmt.Println(re.Name())
	}
}
