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
)

import (
	"github.com/stretchr/testify/assert"
)

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
	got, err := decodeJavaResponse(`customReplyIntegerArray`, ``, false)
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
	got, err := decodeJavaResponse(`customReplyLongArray`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	assert.Equal(t, arr, got)
}
