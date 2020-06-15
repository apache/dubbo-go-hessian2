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

func TestEncBinary(t *testing.T) {
	var (
		v   []byte
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = []byte{}
	e.Encode(v)
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", v, res, err)

	v = []byte{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 'a', 'b', 'c', 'd'}
	e = NewEncoder()
	e.Encode(v)
	t.Logf("encode(%v) = %v\n", v, e.Buffer())
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v, equal:%v\n", v, res, err, bytes.Equal(v, res.([]byte)))
	assertEqual(v, res.([]byte), t)
}

func TestEncBinaryShort(t *testing.T) {
	var (
		v   [1010]byte
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	for i := 0; i < len(v); i++ {
		v[i] = byte(i % 123)
	}

	e = NewEncoder()
	e.Encode(v[:])
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	assertEqual(v[:], res.([]byte), t)
}

func TestEncBinaryChunk(t *testing.T) {
	var (
		v   [65530]byte
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	for i := 0; i < len(v); i++ {
		v[i] = byte(i % 123)
	}

	e = NewEncoder()
	e.Encode(v[:])
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	assertEqual(v[:], res.([]byte), t)
}

func TestBinary(t *testing.T) {
	s0 := ""
	s1 := "0"
	s16 := "0123456789012345"

	s1024 := ""
	for i := 0; i < 16; i++ {
		s1024 += fmt.Sprintf("%02d 456789012345678901234567890123456789012345678901234567890123\n", i)
	}

	s65560 := ""
	for i := 0; i < 1024; i++ {
		s65560 += fmt.Sprintf("%03d 56789012345678901234567890123456789012345678901234567890123\n", i)
	}

	testDecodeFramework(t, "replyBinary_0", []byte(s0))
	testDecodeFramework(t, "replyBinary_1", []byte(s1))
	testDecodeFramework(t, "replyBinary_1023", []byte(s1024[:1023]))
	testDecodeFramework(t, "replyBinary_1024", []byte(s1024))
	testDecodeFramework(t, "replyBinary_15", []byte(s16[:15]))
	testDecodeFramework(t, "replyBinary_16", []byte(s16))
	testDecodeFramework(t, "replyBinary_65536", []byte(s65560[:65536]))
	testDecodeFramework(t, "replyBinary_null", nil)
}

func TestBinaryEncode(t *testing.T) {
	s0 := ""
	s1 := "0"
	s16 := "0123456789012345"

	s1024 := ""
	for i := 0; i < 16; i++ {
		s1024 += fmt.Sprintf("%02d 456789012345678901234567890123456789012345678901234567890123\n", i)
	}

	s65560 := ""
	for i := 0; i < 1024; i++ {
		s65560 += fmt.Sprintf("%03d 56789012345678901234567890123456789012345678901234567890123\n", i)
	}

	testJavaDecode(t, "argBinary_0", []byte(s0))
	testJavaDecode(t, "argBinary_1", []byte(s1))
	testJavaDecode(t, "argBinary_1023", []byte(s1024[:1023]))
	testJavaDecode(t, "argBinary_1024", []byte(s1024))
	testJavaDecode(t, "argBinary_15", []byte(s16[:15]))
	testJavaDecode(t, "argBinary_16", []byte(s16))
	testJavaDecode(t, "argBinary_65536", []byte(s65560[:65536]))
}
