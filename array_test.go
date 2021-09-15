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

func TestBooleanArray(t *testing.T) {
	booleanArray := &BooleanArray{[]bool{true, false}}
	e := &Encoder{}
	a := assert.New(t)

	err := e.Encode(booleanArray)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(booleanArray.Values, decodeValue.(*BooleanArray).Values)
}

func TestIntegerArray(t *testing.T) {
	ia := &IntegerArray{[]int32{1, 2, 3}}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(ia)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(ia.Values, decodeValue.(*IntegerArray).Values)

	// Integer[] that length > 7
	bigIa := &IntegerArray{[]int32{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 3}}
	ne := &Encoder{}
	err = ne.Encode(bigIa)
	a.Nil(err)

	decoder = NewDecoder(e.buffer)
	decodeValue, err = decoder.DecodeValue()
	a.Nil(err)
	a.Equal(ia.Values, decodeValue.(*IntegerArray).Values)
}

func TestByteArray(t *testing.T) {
	ba := &ByteArray{}
	ba.Values = []uint8{1, 2, 3}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(ba)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(ba.Values, decodeValue.(*ByteArray).Values)
}

func TestShortArray(t *testing.T) {
	sa := &ShortArray{}
	sa.Values = []int16{1, 2, 3}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(sa)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(sa.Values, decodeValue.(*ShortArray).Values)
}

func TestLongArray(t *testing.T) {
	la := &LongArray{[]int64{1, 2, 3, 4}}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(la)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(la.Values, decodeValue.(*LongArray).Values)
}

func TestFloatArray(t *testing.T) {
	fa := &FloatArray{[]float32{1, 2, 3, 4}}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(fa)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(fa.Values, decodeValue.(*FloatArray).Values)
}

func TestDoubleArray(t *testing.T) {
	da := &DoubleArray{[]float64{1, 2, 3, 4}}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(da)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(da.Values, decodeValue.(*DoubleArray).Values)
}

func TestCharacterArray(t *testing.T) {
	ca := &CharacterArray{"hello world"}
	a := assert.New(t)

	e := &Encoder{}
	err := e.Encode(ca)
	a.Nil(err)

	decoder := NewDecoder(e.buffer)
	decodeValue, err := decoder.DecodeValue()
	a.Nil(err)
	a.Equal(ca.Values, decodeValue.(*CharacterArray).Values)
}
