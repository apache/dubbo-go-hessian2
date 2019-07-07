package hessian

import (
	"reflect"
	"testing"
)

import (
	big "github.com/dubbogo/gost/math/big"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeDecimal(t *testing.T) {
	var dec big.Decimal
	_ = dec.FromString([]byte("100.256"))
	e := NewEncoder()
	err := e.Encode(dec)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	d := NewDecoder(e.buffer)
	decObj, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(dec.ToString(), decObj.(*big.Decimal).ToString()) {
		t.Errorf("expect: %v, but get: %v", dec, decObj)
	}
}

func TestDecimalFromJava(t *testing.T) {
	var d big.Decimal
	_ = d.FromString([]byte("100.256"))
	d.Value = string(d.ToString())
	doTestDecimal(t, "customReplyTypedFixedDecimal", "100.256")
}

func doTestDecimal(t *testing.T, method, content string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, content, string(r.(*big.Decimal).ToString()))
	})
}
