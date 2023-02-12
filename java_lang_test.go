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
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestDecodePointerSliceSet(t *testing.T) {
	a := []*int16{nil}
	v := reflect.ValueOf(a)
	var b int16 = 1

	typ := reflect.TypeOf(b)

	fmt.Printf("%v %v", v.Index(0).Type(), v.Index(0).Kind())
	switch v.Index(0).Type() {
	case typ:

	}
	v.Index(0).Set(reflect.ValueOf(&b))

}
func TestDecodeJavaSingleShort(t *testing.T) {
	var i int32 = 123
	got, err := decodeJavaResponse(`customReplySingleShort`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, i, got)
}

func TestDecodeJavaShortArray(t *testing.T) {
	var a int16 = 123
	var b int16 = -456

	arr := []*int16{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaShortArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleInteger(t *testing.T) {
	var i int32 = 123
	got, err := decodeJavaResponse(`customReplySingleInteger`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, i, got)
}

func TestDecodeJavaIntegerArray(t *testing.T) {
	var a int32 = 123
	var b int32 = -456

	arr := []*int32{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaIntegerArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleLong(t *testing.T) {
	var i int64 = 12345
	got, err := decodeJavaResponse(`customReplySingleLong`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, i, got)
}

func TestDecodeJavaLongArray(t *testing.T) {
	var a int64 = 12345
	var b int64 = -67890

	arr := []*int64{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaLongArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleBoolean(t *testing.T) {
	var b = true
	got, err := decodeJavaResponse(`customReplySingleBoolean`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, b, got)
}

func TestDecodeJavaBooleanArray(t *testing.T) {
	var a = true
	var b = false

	arr := []*bool{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaBooleanArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleByte(t *testing.T) {
	var b int32 = 'A'
	got, err := decodeJavaResponse(`customReplySingleByte`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, b, got)
}

func TestDecodeJavaByteArray(t *testing.T) {
	var a byte = 'A'
	var b byte = 'C'

	arr := []*byte{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaByteArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleFloat(t *testing.T) {
	var b = 1.23
	got, err := decodeJavaResponse(`customReplySingleFloat`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, b, got)
}

func TestDecodeJavaFloatArray(t *testing.T) {
	var a float32 = 1.23
	var b float32 = 4.56

	arr := []*float32{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaFloatArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleDouble(t *testing.T) {
	var b = 1.23
	got, err := decodeJavaResponse(`customReplySingleDouble`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, b, got)
}

func TestDecodeJavaDoubleArray(t *testing.T) {
	var a = 1.23
	var b = 4.56

	arr := []*float64{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaDoubleArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}

func TestDecodeJavaSingleCharacter(t *testing.T) {
	var b = "A"
	got, err := decodeJavaResponse(`customReplySingleCharacter`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, b, got)
}

func TestDecodeJavaCharacterArray(t *testing.T) {
	var a = 'A'
	var b = 'C'

	arr := []*rune{&a, nil, &b}
	got, err := decodeJavaResponse(`customReplyJavaCharacterArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}
