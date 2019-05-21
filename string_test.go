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
		t.Errorf("Decode() = %v", err)
	}
	// t.Logf("decode(%v) = %v, %v\n", v, res, err)
	assertEqual([]byte(res.(string)), []byte(v), t)
}
