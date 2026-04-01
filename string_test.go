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
	"strings"
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

func TestEncStringChunk(t *testing.T) {
	enc := NewEncoder()
	v := strings.Repeat("我", CHUNK_SIZE-1) + "🤣"
	assert.Nil(t, enc.Encode(v))
	dec := NewDecoder(enc.Buffer())
	s, err := dec.Decode()
	assert.Nil(t, err)
	assert.Equal(t, v, s)
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

func TestDecodeJsonString(t *testing.T) {
	jsonString := `{"params":{"fromAccid":"23495382","msgType":100,"msgId":"148ef1b2-808d-48f2-b268-7a1018a27bdb","attach":"{\"accid\":\"23495382\",\"classRoomFlag\":50685,\"msgId\":\"599645021431398400\",\"msgType\":\"100\",\"nickname\":\"橙子������\"}","roomid":413256699},"url":"https://api.netease.im/nimserver/chatroom/sendMsg.action"}`
	testDecodeFramework(t, "customReplyJsonString", jsonString)
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

func TestStringPrtEncode(t *testing.T) {
	str := "dubbo-go-hessian2"
	testSimpleEncode(t, &str)
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

func TestStringEmoji2(t *testing.T) {
	// see: test_hessian/src/main/java/test/TestString.java
	// see https://github.com/apache/dubbo-go-hessian2/issues/252
	s0 := "❄️🚫🚫🚫🚫 多次自我介绍、任务、动态和"

	testDecodeFramework(t, "customReplyStringEmoji2", s0)
	testJavaDecode(t, "customArgString_emoji2", s0)
}

func TestStringComplex(t *testing.T) {
	// see: test_hessian/src/main/java/test/TestString.java
	s0 := "킐\u0088中国你好!\u0088\u0088\u0088\u0088\u0088\u0088"

	testDecodeFramework(t, "customReplyComplexString", s0)
	testJavaDecode(t, "customArgComplexString", s0)
}

// TestDecodeStringLoneHighSurrogateRealWorld reproduces the real-world string
// "test_go_hessian2\ud83d..." sent by a Java hessian2 encoder.
//
// Java encodes lone surrogates as 3-byte CESU-8: U+D83D → 0xED 0xA0 0xBD.
// The following "..." bytes (0x2E 0x2E 0x2E) decode to c2 = 0x002E, which is
// outside [0xDC00, 0xDFFF].  Before the fix decode2utf8 combined them
// unconditionally, corrupting the suffix; the fix skips the combination.
func TestDecodeStringLoneHighSurrogateRealWorld(t *testing.T) {
	// Byte stream produced by Java hessian2 for "test_go_hessian2\ud83d..."
	serialized := []byte{
		0x14,                                     // BC_STRING_DIRECT, 20 chars
		0x74, 0x65, 0x73, 0x74, 0x5F, 0x67, 0x6F, // test_go
		0x5F, 0x68, 0x65, 0x73, 0x73, 0x69, 0x61, 0x6E, 0x32, // _hessian2
		0xED, 0xA0, 0xBD, // lone high surrogate U+D83D (CESU-8)
		0x2E, 0x2E, 0x2E, // ...
	}

	d := NewDecoder(serialized)
	got, err := d.Decode()
	assert.Nil(t, err)
	assert.Equal(t, "test_go_hessian2\xed\xa0\xbd...", got)
}

func BenchmarkDecodeStringAscii(b *testing.B) {
	runBenchmarkDecodeString(b, "hello world, hello hessian")
}

func BenchmarkDecodeStringUnicode(b *testing.B) {
	runBenchmarkDecodeString(b, "你好, 世界, 你好, hessian")
}

func BenchmarkDecodeStringEmoji(b *testing.B) {
	runBenchmarkDecodeString(b, "❄️🚫🚫🚫🚫 多次自我介绍、任务、动态和")
}

func runBenchmarkDecodeString(b *testing.B, s string) {
	s = strings.Repeat(s, 4096)

	e := NewEncoder()
	_ = e.Encode(s)
	buf := e.buffer

	d := NewDecoder(buf)
	for i := 0; i < b.N; i++ {
		d.Reset(buf)
		_, err := d.Decode()
		if err != nil {
			b.Logf("err: %s", err)
			b.FailNow()
		}
	}
}
