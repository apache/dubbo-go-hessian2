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

package hessian_demo

import (
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

type User struct {
	Name string
	Age  int
}

func (User) JavaClassName() string {
	return "hessian.demo.User"
}

type Location struct {
	Address  string
	Postcode string
}

func (Location) JavaClassName() string {
	return "hessian.demo.Location"
}

func TestLoopEncodeDecode(t *testing.T) {
	u := &User{
		Name: "wongoo",
		Age:  18,
	}
	loc := &Location{
		Address:  "xiamen",
		Postcode: "361000",
	}

	// demo a bytes buffer from client to server
	var bytes []byte
	e := hessian.NewEncoder()

	// -------- encode 1
	bytes = append(bytes, encodeWithFlagAndLength(e, 'D', u)...)

	// -------- encode 2
	bytes = append(bytes, encodeWithFlagAndLength(e, 'S', 12345)...)

	// -------- encode 3
	bytes = append(bytes, encodeWithFlagAndLength(e, 'D', loc)...)

	// -------- encode 4
	bytes = append(bytes, encodeWithFlagAndLength(e, 'D', "hello")...)

	// demo a decoder to decode buffer from client
	d := hessian.NewDecoder(bytes)

	// -------- decode 1
	decodeFlagAndLengthAndData(t, d, u)

	// -------- decode 2
	decodeFlagAndLengthAndData(t, d, 12345)

	// -------- decode 3
	decodeFlagAndLengthAndData(t, d, loc)

	// -------- decode 4
	decodeFlagAndLengthAndData(t, d, "hello")
}

// encode format: [flag][length][binary data]
func encodeWithFlagAndLength(e *hessian.Encoder, flag byte, obj interface{}) []byte {
	var bytes []byte
	bytes = append(bytes, flag)

	e.Clean()
	_ = e.Encode(obj)
	dataBytes := e.Buffer()

	length := len(dataBytes)

	e.Clean()
	_ = e.Encode(length)
	lengthBytes := e.Buffer()

	bytes = append(bytes, lengthBytes...)
	bytes = append(bytes, dataBytes...)

	return bytes
}

func decodeFlagAndLengthAndData(t *testing.T, d *hessian.Decoder, expect interface{}) {
	// decode flag
	d.Clean()
	flag, _ := d.ReadByte()

	// decode length
	d.Clean()
	lengthObj, _ := d.Decode()
	length := lengthObj.(int64)

	// skip data when flag='S'
	if flag == 'S' {
		_, _ = d.Discard(int(length))
		return
	}

	// decode data
	d.Clean()
	res, _ := d.Decode()

	// check
	assert.True(t, reflect.DeepEqual(expect, res))
}
