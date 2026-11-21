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
	"bytes"
	"fmt"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

var assertEqual = func(want, got []byte, t *testing.T) {
	if !bytes.Equal(want, got) {
		t.Fatalf("want %v , got %v", want, got)
	}
}

func encodeTarget(target interface{}) ([]byte, error) {
	e := NewEncoder()
	err := e.Encode(target)
	if err != nil {
		return nil, err
	}
	return e.Buffer(), nil
}

func javaDecodeValidate(t *testing.T, method string, target interface{}) (string, error) {
	b, err := encodeTarget(target)
	if err != nil {
		return "", err
	}

	genHessianJar()
	cmd, err := newJavaCommand("-jar", hessianJar, method)
	if err != nil {
		t.Logf("resolve java error: %v", err)
		return "", err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Logf("call java error: %v", err)
		return "", err
	}

	go func() {
		_, werr := stdin.Write(b)
		if werr != nil {
			panic(fmt.Errorf("write to stdin error: %v", werr))
		}
		_ = stdin.Close()
	}()

	out, err := cmd.Output()
	if err != nil {
		t.Logf("get java result error: %v", err)
		return "", err
	}
	return string(out), nil
}

func testJavaDecode(t *testing.T, method string, target interface{}) {
	result, err := javaDecodeValidate(t, method, target)
	if err != nil {
		t.Errorf("%s: encode fail with error: %v", method, err)
		return
	}

	if result != "true" {
		t.Errorf("%s: encode %v to bytes wrongly, result: %s", method, target, result)
	}
}

func testSimpleEncode(t *testing.T, v interface{}) {
	e := NewEncoder()
	err := e.Encode(v)
	assert.Nil(t, err)
}

type BenchData struct {
	name string
}

func (b *BenchData) JavaClassName() string {
	return "test.bench.BenchData"
}

// Benchmark_Struct_Encode  	 2231869	       506.8 ns/op	     560 B/op	       7 allocs/op
func Benchmark_Struct_Encode(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewEncoder().Encode(BenchData{})
	}
}

// Benchmark_Pointer_Encode   	 2565778	       476.1 ns/op	     560 B/op	       7 allocs/op
func Benchmark_Pointer_Encode(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewEncoder().Encode(&BenchData{})
	}
}

// Benchmark_Struct_Encode_8   	 2307214	       519.4 ns/op	     560 B/op	       7 allocs/op
func Benchmark_Struct_Encode_8(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewEncoder().Encode(BenchData{})
		}
	})
}

// Benchmark_Pointer_Encode_8   	 2460842	       476.7 ns/op	     560 B/op	       7 allocs/op
func Benchmark_Pointer_Encode_8(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewEncoder().Encode(&BenchData{})
		}
	})
}

type (
	namedString  string
	namedInt     int
	namedInt8    int8
	namedInt16   int16
	namedInt64   int64
	namedUint    uint
	namedUint8   uint8
	namedUint16  uint16
	namedUint32  uint32
	namedUint64  uint64
	namedFloat32 float32
	namedFloat64 float64
	namedBool    bool
)

// https://github.com/apache/dubbo-go-hessian2/issues/404
func TestEncodeNamedBasicTypes(t *testing.T) {
	s := namedString("ok")
	f := namedBool(true)
	b := true
	pb := &b

	cases := []struct {
		named interface{}
		plain interface{}
	}{
		{namedString("ok"), "ok"},
		{&s, "ok"},
		{(*namedString)(nil), nil},
		{namedInt(-7), -7},
		{namedInt8(-8), int8(-8)},
		{namedInt16(300), int16(300)},
		{namedInt64(1 << 40), int64(1 << 40)},
		{namedUint(7), uint(7)},
		{namedUint8(200), uint8(200)},
		{namedUint16(60000), uint16(60000)},
		{namedUint32(1 << 31), uint32(1 << 31)},
		{namedUint64(1 << 40), uint64(1 << 40)},
		{namedFloat32(1.5), float32(1.5)},
		{namedFloat64(2.25), 2.25},
		{namedBool(true), true},
		{&f, true},
		{(*namedBool)(nil), nil},
		{&pb, true},
	}
	for _, c := range cases {
		want, err := encodeTarget(c.plain)
		assert.Nil(t, err)
		got, err := encodeTarget(c.named)
		assert.Nil(t, err)
		assert.Equal(t, want, got, "%T", c.named)
	}
}

type NamedTypes struct {
	Status    namedString
	Code      namedInt64
	Flag      namedBool
	Ratio     namedFloat64
	PtrStatus *namedString
	PtrCode   *namedInt64
	PtrFlag   *namedBool
	Name      string
}

func (*NamedTypes) JavaClassName() string {
	return "test.model.NamedTypes"
}

func TestNamedTypesStruct(t *testing.T) {
	RegisterPOJO(&NamedTypes{})
	status, code, flag := namedString("PENDING"), namedInt64(7), namedBool(false)
	want := &NamedTypes{
		Status:    "ACTIVE",
		Code:      42,
		Flag:      true,
		Ratio:     1.5,
		PtrStatus: &status,
		PtrCode:   &code,
		PtrFlag:   &flag,
		Name:      "bob",
	}

	b, err := encodeTarget(want)
	assert.Nil(t, err)
	assert.Equal(t, want, mustDecodeObject(t, b))

	testJavaDecode(t, "customArgNamedTypes", want)
	testDecodeFramework(t, "customReplyNamedTypes", want)
}
