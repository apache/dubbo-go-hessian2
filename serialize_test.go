/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"fmt"
	"reflect"
	"testing"

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

	if !reflect.DeepEqual(dec.String(), decObj.(*big.Decimal).String()) {
		t.Errorf("expect: %v, but get: %v", dec, decObj)
	}
}

func TestDecimalGoDecode(t *testing.T) {
	var d big.Decimal
	_ = d.FromString("100.256")
	d.Value = d.String()
	doTestStringer(t, "customReplyTypedFixedDecimal", "100.256")
}

func TestDecimalJavaDecode(t *testing.T) {
	var d big.Decimal
	_ = d.FromString("100.256")
	d.Value = d.String()
	testJavaDecode(t, "customArgTypedFixed_Decimal", d)
}

func TestEncodeDecodeInteger(t *testing.T) {
	var dec big.Integer
	_ = dec.FromString("100256")
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

	if !reflect.DeepEqual(dec.String(), decObj.(*big.Integer).String()) {
		t.Errorf("expect: %v, but get: %v", dec, decObj)
	}
}

func TestIntegerGoDecode(t *testing.T) {
	var d big.Integer
	_ = d.FromString("100256")
	d.Value = d.String()
	doTestStringer(t, "customReplyTypedFixedInteger", "100256")
}

func TestIntegerJavaDecode(t *testing.T) {
	var d big.Integer
	_ = d.FromString("100256")
	d.Value = d.String()
	testJavaDecode(t, "customArgTypedFixed_Integer", d)
}

func doTestStringer(t *testing.T, method, content string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, content, r.(fmt.Stringer).String())
	})
}
