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

func TestLong(t *testing.T) {
	testDecodeFramework(t, "replyLong_0", int64(0))
	testDecodeFramework(t, "replyLong_0x10", int64(0x10))
	testDecodeFramework(t, "replyLong_0x3ffff", int64(0x3ffff))
	testDecodeFramework(t, "replyLong_0x40000", int64(0x40000))
	testDecodeFramework(t, "replyLong_0x7ff", int64(0x7ff))
	testDecodeFramework(t, "replyLong_0x7fffffff", int64(0x7fffffff))
	testDecodeFramework(t, "replyLong_0x800", int64(0x800))
	testDecodeFramework(t, "replyLong_1", int64(1))
	testDecodeFramework(t, "replyLong_15", int64(15))
	testDecodeFramework(t, "replyLong_m0x40000", int64(-0x40000))
	testDecodeFramework(t, "replyLong_m0x40001", int64(-0x40001))
	testDecodeFramework(t, "replyLong_m0x800", int64(-0x800))
	testDecodeFramework(t, "replyLong_m0x80000000", int64(-0x80000000))
	testDecodeFramework(t, "replyLong_m0x80000001", int64(-0x80000001))
	testDecodeFramework(t, "replyLong_m0x801", int64(-0x801))
	testDecodeFramework(t, "replyLong_m8", int64(-8))
	testDecodeFramework(t, "replyLong_m9", int64(-9))
}

func TestLongEncode(t *testing.T) {
	testJavaDecode(t, "argLong_0", int64(0))
	testJavaDecode(t, "argLong_0x10", int64(0x10))
	testJavaDecode(t, "argLong_0x3ffff", int64(0x3ffff))
	testJavaDecode(t, "argLong_0x40000", int64(0x40000))
	testJavaDecode(t, "argLong_0x7ff", int64(0x7ff))
	testJavaDecode(t, "argLong_0x7fffffff", int64(0x7fffffff))
	testJavaDecode(t, "argLong_0x800", int64(0x800))
	testJavaDecode(t, "argLong_1", int64(1))
	testJavaDecode(t, "argLong_15", int64(15))
	testJavaDecode(t, "argLong_m0x40000", int64(-0x40000))
	testJavaDecode(t, "argLong_m0x40001", int64(-0x40001))
	testJavaDecode(t, "argLong_m0x800", int64(-0x800))
	testJavaDecode(t, "argLong_m0x80000000", int64(-0x80000000))
	testJavaDecode(t, "argLong_m0x80000001", int64(-0x80000001))
	testJavaDecode(t, "argLong_m0x801", int64(-0x801))
	testJavaDecode(t, "argLong_m8", int64(-8))
	testJavaDecode(t, "argLong_m9", int64(-9))
}
