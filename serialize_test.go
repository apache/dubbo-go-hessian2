// Copyright 2016-2019 aliiohs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	_ = dec.FromString("100.256")
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

func TestDecimalGoDecode(t *testing.T) {
	var d big.Decimal
	_ = d.FromString("100.256")
	d.Value = d.String()
	doTestDecimal(t, "customReplyTypedFixedDecimal", "100.256")
}

func TestDecimalJavaDecode(t *testing.T) {
	var d big.Decimal
	_ = d.FromString("100.256")
	d.Value = d.String()
	testJavaDecode(t, "customArgTypedFixedList_Decimal", d)
}

func doTestDecimal(t *testing.T, method, content string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, content, string(r.(*big.Decimal).ToString()))
	})
}
