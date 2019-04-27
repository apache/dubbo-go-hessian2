// Copyright (c) 2016 ~ 2019, Alex Stocks.
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

func TestEncInt32Len1B(t *testing.T) {
	var (
		v   int32
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	v = 0xe6
	// var v int32 = 0xf016
	e = NewEncoder()
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func TestEncInt32Len2B(t *testing.T) {
	var (
		v   int32
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	v = 0xf016
	e = NewEncoder()
	e.Encode(v)
	if len(e.buffer) == 0 {
		t.Fail()
	}
	t.Logf("%#v\n", e.buffer)
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%#x) = %#x, %v\n", v, res, err)
}

func TestEncInt32Len4B(t *testing.T) {
	var (
		v   int32
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 0x20161024
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func testIntFramework(t *testing.T, method string, expected int32) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(int32)
	if !ok {
		t.Errorf("%s: %v is not int", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestInt(t *testing.T) {
	testIntFramework(t, "replyInt_0", 0)
	testIntFramework(t, "replyInt_0x30", 0x30)
	testIntFramework(t, "replyInt_0x3ffff", 0x3ffff)
	testIntFramework(t, "replyInt_0x40000", 0x40000)
	testIntFramework(t, "replyInt_0x7ff", 0x7ff)
	testIntFramework(t, "replyInt_0x7fffffff", 0x7fffffff)
	testIntFramework(t, "replyInt_0x800", 0x800)
	testIntFramework(t, "replyInt_1", 1)
	testIntFramework(t, "replyInt_47", 47)
	testIntFramework(t, "replyInt_m0x40000", -0x40000)
	testIntFramework(t, "replyInt_m0x40001", -0x40001)
	testIntFramework(t, "replyInt_m0x800", -0x800)
	testIntFramework(t, "replyInt_m0x80000000", -0x80000000)
	testIntFramework(t, "replyInt_m0x801", -0x801)
	testIntFramework(t, "replyInt_m16", -16)
	testIntFramework(t, "replyInt_m17", -17)
}
