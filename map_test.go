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
	"testing"
)

func TestEncUntypedMap(t *testing.T) {
	var (
		m   map[interface{}]interface{}
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	m = make(map[interface{}]interface{})
	m["hello"] = "world"
	m[100] = "100"
	m[100.1010] = 101910
	m[true] = true
	m[false] = true
	e.Encode(m)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", m, res, err)
}

func TestEncTypedMap(t *testing.T) {
	var (
		m   map[int]string
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	m = make(map[int]string)
	m[0] = "hello"
	m[1] = "golang"
	m[2] = "world"
	e.Encode(m)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", m, res, err)
}
