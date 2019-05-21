// Copyright 2016-2019 Alex Stocks, Xinge Gao
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
	"testing"
)

func TestEncDouble(t *testing.T) {
	var (
		v   float64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 2016.1024
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func testDoubleFramework(t *testing.T, method string, expected float64) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %+v", method, e)
		return
	}

	v, ok := r.(float64)
	if !ok {
		t.Errorf("%s: %v is not double", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestDouble(t *testing.T) {
	testDoubleFramework(t, "replyDouble_0_0", 0.0)
	testDoubleFramework(t, "replyDouble_0_001", 0.001)
	testDoubleFramework(t, "replyDouble_1_0", 1.0)
	testDoubleFramework(t, "replyDouble_127_0", 127.0)
	testDoubleFramework(t, "replyDouble_128_0", 128.0)
	testDoubleFramework(t, "replyDouble_2_0", 2.0)
	testDoubleFramework(t, "replyDouble_3_14159", 3.14159)
	testDoubleFramework(t, "replyDouble_32767_0", 32767.0)
	testDoubleFramework(t, "replyDouble_65_536", 65.536)
	testDoubleFramework(t, "replyDouble_m0_001", -0.001)
	testDoubleFramework(t, "replyDouble_m128_0", -128.0)
	testDoubleFramework(t, "replyDouble_m129_0", -129.0)
	testDoubleFramework(t, "replyDouble_m32768_0", -32768.0)
}
