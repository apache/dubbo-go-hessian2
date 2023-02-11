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

	"github.com/stretchr/testify/assert"
)

func TestBooleanArray(t *testing.T) {
	var x, y = true, false
	booleanArray := []*bool{&x, &y}
	e := NewEncoder()

	err := e.Encode(booleanArray)
	assert.Nil(t, err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	assert.Nil(t, err)
	assert.Equal(t, booleanArray, decodeValue)
}

func TestIntegerArray(t *testing.T) {
	var a, b, c int32 = 1, 2, 3
	ia := []*int32{&a, &b, &c}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(ia)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(ia, decodeValue)

	// Integer[] that length > 7
	//bigIa := &IntegerArray{[]int32{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}}
	bigIa := []*int32{&a, &b, &c, &a, &b, &c, &a, &b, &c, &a, &b, &c, &a, &b, &c, &a, &b, &c}
	ne := NewEncoder()
	err = ne.Encode(bigIa)
	tt.Nil(err)

	decoder = NewDecoder(e.buffer)
	decodeValue, err = decoder.Decode()
	tt.Nil(err)
	tt.Equal(ia, decodeValue)
}

func TestByteArray(t *testing.T) {
	var a, b, c uint8 = 1, 2, 3
	ba := []*byte{&a, &b, &c}

	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(ba)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(ba, decodeValue)
}

func TestShortArray(t *testing.T) {
	var a, b, c int16 = 1, 2, 3
	sa := []*int16{&a, &b, &c}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(sa)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(sa, decodeValue)
}

func TestLongArray(t *testing.T) {
	var a, b, c, d int64 = 1, 2, 3, 4
	la := []*int64{&a, &b, &c, &d}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(la)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(la, decodeValue)
}

func TestFloatArray(t *testing.T) {
	var a, b, c, d float32 = 1, 2, 3, 4
	fa := []*float32{&a, &b, &c, &d}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(fa)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(fa, decodeValue)
}

func TestDoubleArray(t *testing.T) {
	var a, b, c, d float64 = 1, 2, 3, 4
	da := []*float64{&a, &b, &c, &d}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(da)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)
	tt.Equal(da, decodeValue)
}

func TestCharacterArray(t *testing.T) {
	var r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11 Rune = 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'
	ca := []*Rune{&r1, &r2, &r3, &r4, &r5, &r6, &r7, &r8, &r9, &r10, &r11}
	tt := assert.New(t)

	e := NewEncoder()
	err := e.Encode(ca)
	tt.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.Decode()
	tt.Nil(err)

	expected := []*int32{(*int32)(&r1), (*int32)(&r2), (*int32)(&r3), (*int32)(&r4), (*int32)(&r5), (*int32)(&r6), (*int32)(&r7), (*int32)(&r8), (*int32)(&r9), (*int32)(&r10), (*int32)(&r11)}
	tt.Equal(expected, decodeValue)
}
