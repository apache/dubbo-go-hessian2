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

func TestNullTypedScalarPointers(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{name: "int32", value: (*int32)(nil)},
		{name: "int64", value: (*int64)(nil)},
		{name: "bool", value: (*bool)(nil)},
		{name: "string", value: (*string)(nil)},
		{name: "float64", value: (*float64)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEncoder()
			if err := e.Encode(tt.value); err != nil {
				t.Fatalf("encode typed nil %s pointer: %v", tt.name, err)
			}
			if e.Buffer() == nil {
				t.Fatal("expected encoder buffer")
			}
			assertEqual([]byte("N"), e.buffer, t)
		})
	}
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
	Int64 *int64
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
	assertEqual([]byte("NNNNNN"), e.buffer[len(e.buffer)-6:], t)
}

// Int64PtrFieldStruct verifies that a *int64 struct field is encoded as a
// hessian long, so that a java consumer can deserialize it into a
// java.lang.Long field, see apache/dubbo-go#2410.
type Int64PtrFieldStruct struct {
	Total *int64
}

func (*Int64PtrFieldStruct) JavaClassName() string {
	return "Int64PtrFieldStruct"
}

func TestInt64PtrFieldStructEncode(t *testing.T) {
	total := int64(12345)
	e := NewEncoder()
	if err := e.Encode(&Int64PtrFieldStruct{Total: &total}); err != nil {
		t.Fatalf("encode Int64PtrFieldStruct: %v", err)
	}
	// 0x3c 0x30 0x39 is the hessian short-form long encoding of 12345
	assertEqual([]byte{0x3c, 0x30, 0x39}, e.buffer[len(e.buffer)-3:], t)

	d := NewDecoder(e.Buffer())
	v, err := d.Decode()
	if err != nil {
		t.Fatalf("decode Int64PtrFieldStruct: %v", err)
	}
	got, ok := v.(*Int64PtrFieldStruct)
	if !ok {
		t.Fatalf("decode Int64PtrFieldStruct: unexpected type %T", v)
	}
	if got.Total == nil || *got.Total != total {
		t.Fatalf("decode Int64PtrFieldStruct: unexpected Total: %+v", got.Total)
	}
}
