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

func TestEncNull(t *testing.T) {
	e := NewEncoder()
	e.Encode(nil)
	if e.Buffer() == nil {
		t.Fail()
	}
	t.Logf("nil enc result:%s\n", string(e.buffer))
}

func TestNull(t *testing.T) {
	testDecodeFramework(t, "replyNull", nil)
}

func TestNullEncode(t *testing.T) {
	testJavaDecode(t, "argNull", nil)
}

func TestNullIntPtr(t *testing.T) {
	e := NewEncoder()
	var null *int = nil
	e.Encode(null)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("N"), e.buffer, t)
}

func TestNullBoolPtr(t *testing.T) {
	e := NewEncoder()
	var null *bool = nil
	e.Encode(null)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("N"), e.buffer, t)
}

func TestNullInt32Ptr(t *testing.T) {
	e := NewEncoder()
	var null *int32 = nil
	e.Encode(null)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("N"), e.buffer, t)
}

func TestNullSlice(t *testing.T) {
	e := NewEncoder()
	var null []int32 = nil
	e.Encode(null)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("N"), e.buffer, t)
}

func TestNullMap(t *testing.T) {
	e := NewEncoder()
	var null map[bool]int32 = nil
	e.Encode(null)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("N"), e.buffer, t)
}

type NullFieldStruct struct {
	Int   *int
	Bool  *bool
	Int32 *int32
	Slice []int32
	Map   map[bool]int32
}

func (*NullFieldStruct) JavaClassName() string {
	return "NullFieldStruct"
}

func TestNullFieldStruct(t *testing.T) {
	e := NewEncoder()
	req := &NullFieldStruct{}
	e.Encode(req)
	if e.Buffer() == nil {
		t.Fail()
	}
	assertEqual([]byte("NNNNN"), e.buffer[len(e.buffer)-5:], t)
}
