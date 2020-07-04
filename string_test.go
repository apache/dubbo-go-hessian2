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
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestEncString(t *testing.T) {
	var (
		v   string
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = "hello"
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func TestEncShortRune(t *testing.T) {
	var (
		v   string
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = "我化尘埃飞扬，追寻赤裸逆翔"
	err = e.Encode(v)
	assert.Nil(t, err)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	assert.Nil(t, err)
	if err != nil {
		t.Logf("err:%s", err.Error())
	}
	assert.Equal(t, v, res)
}

func TestEncRune(t *testing.T) {
	var (
		v   string
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()

	v = "我化尘埃飞扬，追寻赤裸逆翔,奔去七月刑场，时间烧灼滚烫,回忆撕毁臆想，路上行走匆忙,难能可贵世上，散播留香磁场,我欲乘风破浪，踏遍黄沙海洋,与其误会一场，也要不负勇往,我愿你是个谎，从未出现南墙,笑是神的伪装，笑是强忍的伤,我想你就站在，站在大漠边疆,我化尘埃飞扬，追寻赤裸逆翔," +
		"奔去七月刑场，时间烧灼滚烫,回忆撕毁臆想，路上行走匆忙,难能可贵世上，散播留香磁场,我欲乘风破浪，踏遍黄沙海洋,与其误会一场，也要不负勇往,我愿你是个谎，从未出现南墙,笑是神的伪装，笑是强忍的伤,我想你就站在，站在大漠边疆."

	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v

	t.Logf("TestEncRune vlen: %d", len(v))

	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}

	assertEqual([]byte(res.(string)), []byte(v), t)
}

func TestString(t *testing.T) {
	s0 := ""
	s1 := "0"
	s32 := "01234567890123456789012345678901"

	s1024 := ""
	for i := 0; i < 16; i++ {
		s1024 += fmt.Sprintf("%02d 456789012345678901234567890123456789012345678901234567890123\n", i)
	}

	s65560 := ""
	for i := 0; i < 1024; i++ {
		s65560 += fmt.Sprintf("%03d 56789012345678901234567890123456789012345678901234567890123\n", i)
	}

	testDecodeFramework(t, "replyString_0", s0)
	testDecodeFramework(t, "replyString_1", s1)
	testDecodeFramework(t, "replyString_1023", s1024[:1023])
	testDecodeFramework(t, "replyString_1024", s1024)
	testDecodeFramework(t, "replyString_31", s32[:31])
	testDecodeFramework(t, "replyString_32", s32)
	testDecodeFramework(t, "replyString_65536", s65560[:65536])
	testDecodeFramework(t, "replyString_null", nil)
}

func TestStringEncode(t *testing.T) {
	s0 := ""
	s1 := "0"
	s32 := "01234567890123456789012345678901"

	s1024 := ""
	for i := 0; i < 16; i++ {
		s1024 += fmt.Sprintf("%02d 456789012345678901234567890123456789012345678901234567890123\n", i)
	}

	s65560 := ""
	for i := 0; i < 1024; i++ {
		s65560 += fmt.Sprintf("%03d 56789012345678901234567890123456789012345678901234567890123\n", i)
	}

	testJavaDecode(t, "argString_0", s0)
	testJavaDecode(t, "argString_1", s1)
	testJavaDecode(t, "argString_1023", s1024[:1023])
	testJavaDecode(t, "argString_1024", s1024)
	testJavaDecode(t, "argString_31", s32[:31])
	testJavaDecode(t, "argString_32", s32)
	testJavaDecode(t, "argString_65536", s65560[:65536])
}

var decodePool = &sync.Pool{
	New: func() interface{} {
		return NewCheapDecoderWithSkip([]byte{})
	},
}

func TestStringWithPool(t *testing.T) {
	e := NewEncoder()
	e.Encode(testString)
	buf := e.buffer

	for i := 0; i < 3; i++ {
		d := decodePool.Get().(*Decoder)
		d.Reset(buf)

		v, err := d.Decode()
		if err != nil {
			t.Errorf("err:%s", err.Error())
		}
		if v != testString {
			t.Errorf("excpect decode %v, actual %v", testString, v)
		}

		decodePool.Put(d)
	}

}

func TestStringEmoji(t *testing.T) {
	// see: test_hessian/src/main/java/test/TestString.java
	s0 := "emoji🤣"
	s0 += ",max" + string(rune(0x10FFFF))

	testDecodeFramework(t, "customReplyStringEmoji", s0)
	testJavaDecode(t, "customArgString_emoji", s0)
}

func TestStringComplex(t *testing.T) {
	// see: test_hessian/src/main/java/test/TestString.java
	s0 := "킐\u0088中国你好!\u0088\u0088\u0088\u0088\u0088\u0088"

	testDecodeFramework(t, "customReplyComplexString", s0)
	testJavaDecode(t, "customArgComplexString", s0)
}
