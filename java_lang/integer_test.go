package java_lang

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerReflect(t *testing.T) {
	a := Integer(1234)
	dest := reflect.ValueOf(&a)
	dest.Elem().SetInt(5678)
	assert.Equal(t, Integer(5678), a)

	sli := []*Integer{nil, nil, nil}
	dest = reflect.ValueOf(sli)

	b := Integer(5678)
	dest.Index(0).Set(reflect.ValueOf(&b))
	dest.Index(1).Set(reflect.ValueOf(&b))
	dest.Index(2).Set(reflect.ValueOf(&b))

	assert.Equal(t, []*Integer{&b, &b, &b}, sli)
}
