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
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	t.Logf("decode(%#x) = %#x, %v\n", v, res, err)
}

func TestEncInt32ForAlias(t *testing.T) {
	type JavaEnum int32
	var (
		v   JavaEnum
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
	assert.Nil(t, err)
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
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
	assert.Nil(t, err)
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func TestInt(t *testing.T) {
	testDecodeFramework(t, "replyInt_0", int32(0))
	testDecodeFramework(t, "replyInt_0x30", int32(0x30))
	testDecodeFramework(t, "replyInt_0x3ffff", int32(0x3ffff))
	testDecodeFramework(t, "replyInt_0x40000", int32(0x40000))
	testDecodeFramework(t, "replyInt_0x7ff", int32(0x7ff))
	testDecodeFramework(t, "replyInt_0x7fffffff", int32(0x7fffffff))
	testDecodeFramework(t, "replyInt_0x800", int32(0x800))
	testDecodeFramework(t, "replyInt_1", int32(1))
	testDecodeFramework(t, "replyInt_47", int32(47))
	testDecodeFramework(t, "replyInt_m0x40000", int32(-0x40000))
	testDecodeFramework(t, "replyInt_m0x40001", int32(-0x40001))
	testDecodeFramework(t, "replyInt_m0x800", int32(-0x800))
	testDecodeFramework(t, "replyInt_m0x80000000", int32(-0x80000000))
	testDecodeFramework(t, "replyInt_m0x801", int32(-0x801))
	testDecodeFramework(t, "replyInt_m16", int32(-16))
	testDecodeFramework(t, "replyInt_m17", int32(-17))
}

func TestIntEncode(t *testing.T) {
	testJavaDecode(t, "argInt_0", int32(0))
	testJavaDecode(t, "argInt_0x30", int32(0x30))
	testJavaDecode(t, "argInt_0x3ffff", int32(0x3ffff))
	testJavaDecode(t, "argInt_0x40000", int32(0x40000))
	testJavaDecode(t, "argInt_0x7ff", int32(0x7ff))
	testJavaDecode(t, "argInt_0x7fffffff", int32(0x7fffffff))
	testJavaDecode(t, "argInt_0x800", int32(0x800))
	testJavaDecode(t, "argInt_1", int32(1))
	testJavaDecode(t, "argInt_47", int32(47))
	testJavaDecode(t, "argInt_m0x40000", int32(-0x40000))
	testJavaDecode(t, "argInt_m0x40001", int32(-0x40001))
	testJavaDecode(t, "argInt_m0x800", int32(-0x800))
	testJavaDecode(t, "argInt_m0x80000000", int32(-0x80000000))
	testJavaDecode(t, "argInt_m0x801", int32(-0x801))
	testJavaDecode(t, "argInt_m16", int32(-16))
	testJavaDecode(t, "argInt_m17", int32(-17))
}

func TestReflectIntEncode(t *testing.T) {
	a1 := int32(0)
	a2 := int32(0x30)
	a3 := int32(0x3ffff)
	a4 := int32(0x40000)
	a5 := int32(0x7ff)
	a6 := int32(0x7fffffff)
	a7 := int32(0x800)
	a8 := int32(1)
	a9 := int32(47)
	a10 := int32(-0x40000)
	a11 := int32(-0x40001)
	a12 := int32(-0x800)
	a13 := int32(-0x80000000)
	a14 := int32(-0x801)
	a15 := int32(-16)
	a16 := int32(-17)
	testJavaDecode(t, "argInt_0", &a1)
	testJavaDecode(t, "argInt_0x30", &a2)
	testJavaDecode(t, "argInt_0x3ffff", &a3)
	testJavaDecode(t, "argInt_0x40000", &a4)
	testJavaDecode(t, "argInt_0x7ff", &a5)
	testJavaDecode(t, "argInt_0x7fffffff", &a6)
	testJavaDecode(t, "argInt_0x800", &a7)
	testJavaDecode(t, "argInt_1", &a8)
	testJavaDecode(t, "argInt_47", &a9)
	testJavaDecode(t, "argInt_m0x40000", &a10)
	testJavaDecode(t, "argInt_m0x40001", &a11)
	testJavaDecode(t, "argInt_m0x800", &a12)
	testJavaDecode(t, "argInt_m0x80000000", &a13)
	testJavaDecode(t, "argInt_m0x801", &a14)
	testJavaDecode(t, "argInt_m16", &a15)
	testJavaDecode(t, "argInt_m17", &a16)
}
