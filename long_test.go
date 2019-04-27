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

func TestEncInt64Len1BDirect(t *testing.T) {
	var (
		v   int64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 0x1
	e.Encode(int64(v))
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(int64(%#x)) = %#x, %v\n", v, res, err)
}

func TestEncInt64Len1B(t *testing.T) {
	var (
		v   int64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 0xf6
	e.Encode(int64(v))
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(int64(%#x)) = %#x, %v\n", v, res, err)
}

func TestEncInt64Len2B(t *testing.T) {
	var (
		v   int64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 0x2016
	e.Encode(int64(v))
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(int64(%#x)) = %#x, %v\n", v, res, err)
}

func TestEncInt64Len3B(t *testing.T) {
	var (
		v   int64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 101910 // 0x18e16
	e.Encode(int64(v))
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(int64(%#x)) = %#x, %v\n", v, res, err)
}

func TestEncInt64Len8B(t *testing.T) {
	var (
		v   int64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 0x20161024114530
	e.Encode(int64(v))
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(int64(%#x)) = %#x, %v\n", v, res, err)
}

func testLongFramework(t *testing.T, method string, expected int64) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(int64)
	if !ok {
		t.Errorf("%s: %v is not long", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestLong(t *testing.T) {
	testLongFramework(t, "replyLong_0", 0)
	testLongFramework(t, "replyLong_0x10", 0x10)
	testLongFramework(t, "replyLong_0x3ffff", 0x3ffff)
	testLongFramework(t, "replyLong_0x40000", 0x40000)
	testLongFramework(t, "replyLong_0x7ff", 0x7ff)
	testLongFramework(t, "replyLong_0x7fffffff", 0x7fffffff)
	testLongFramework(t, "replyLong_0x800", 0x800)
	testLongFramework(t, "replyLong_1", 1)
	testLongFramework(t, "replyLong_15", 15)
	testLongFramework(t, "replyLong_m0x40000", -0x40000)
	testLongFramework(t, "replyLong_m0x40001", -0x40001)
	testLongFramework(t, "replyLong_m0x800", -0x800)
	testLongFramework(t, "replyLong_m0x80000000", -0x80000000)
	testLongFramework(t, "replyLong_m0x80000001", -0x80000001)
	testLongFramework(t, "replyLong_m0x801", -0x801)
	testLongFramework(t, "replyLong_m8", -8)
	testLongFramework(t, "replyLong_m9", -9)
}
