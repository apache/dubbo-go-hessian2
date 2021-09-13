package hessian

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	var intArr = []int32{1, 2, 3}
	integerArray := IntegerArray{intArr}
	jcs := JavaCollectionSerializer{}
	e := &Encoder{}
	err := jcs.EncObject(e, integerArray)
	if err != nil {
		return
	}
	fmt.Printf("%v\n", e.buffer)
}
