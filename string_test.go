// Copyright 2016-2019 Alex Stocks
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
	"fmt"
	"testing"
)
import (
	"github.com/stretchr/testify/assert"
)

func init() {
	RegisterPOJO(&Demo{})
}

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
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
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
	v = "我化尘埃飞扬，追寻赤裸逆翔, 奔去七月刑场，时间烧灼滚烫, 回忆撕毁臆想，路上行走匆忙, 难能可贵世上，散播留香磁场, 我欲乘风破浪，踏遍黄沙海洋, 与其误会一场，也要不负勇往, 我愿你是个谎，从未出现南墙, 笑是神的伪装，笑是强忍的伤, 我想你就站在，站在大漠边疆, 我化尘埃飞扬，追寻赤裸逆翔," +
		" 奔去七月刑场，时间烧灼滚烫, 回忆撕毁臆想，路上行走匆忙, 难能可贵世上，散播留香磁场, 我欲乘风破浪，踏遍黄沙海洋, 与其误会一场，也要不负勇往, 我愿你是个谎，从未出现南墙, 笑是神的伪装，笑是强忍的伤, 我想你就站在，站在大漠边疆."
	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v
	v = v + v + v + v + v
	fmt.Printf("vlen:%d\n", len(v))
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

type Demo struct {
	Name    *string
	Address string
	Age     int
}

func (Demo) JavaClassName() string {
	return "test.model.StringDemo"
}

func TestStringNilDecode(t *testing.T) {

	doTestStringNull(t, "customReplyTypedFixedStringNull")
}

func TestNilPointerForStringEncodeAndDecode(t *testing.T) {
	var s1 *string = nil
	demo := Demo{
		Name:    s1,
		Address: "",
		Age:     12,
	}
	e := NewEncoder()
	_ = e.Encode(demo)
	d := NewDecoder(e.buffer)
	result, _ := d.Decode()
	assert.Equal(t, true, result.(*Demo).Name == nil)
	assert.Equal(t, true, result.(*Demo).Address == "")
}

func doTestStringNull(t *testing.T, method string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, true, r.(*Demo).Name == nil)
		assert.Equal(t, true, r.(*Demo).Address == "")
	})
}
