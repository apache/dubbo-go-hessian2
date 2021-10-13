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

func TestEncDouble(t *testing.T) {
	var (
		v   float64
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 2016.1024
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%v) = %v, %v\n", v, res, err)
}

func TestIssue181(t *testing.T) {
	var (
		v   float32
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = 99.8
	e.Encode(v)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	// res would be '99.800003' without patches in PR #196
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	f := res.(float64)
	if float32(f) != v {
		t.Errorf("decode(%v) = %v, %v\n", v, res, err)
		return
	}
	t.Logf("decode(%v) = %v\n", v, res)
}

func TestDouble(t *testing.T) {
	testDecodeFramework(t, "replyDouble_0_0", 0.0)
	testDecodeFramework(t, "replyDouble_0_001", 0.001)
	testDecodeFramework(t, "replyDouble_1_0", 1.0)
	testDecodeFramework(t, "replyDouble_127_0", 127.0)
	testDecodeFramework(t, "replyDouble_128_0", 128.0)
	testDecodeFramework(t, "replyDouble_2_0", 2.0)
	testDecodeFramework(t, "replyDouble_3_14159", 3.14159)
	testDecodeFramework(t, "replyDouble_32767_0", 32767.0)
	testDecodeFramework(t, "replyDouble_65_536", 65.536)
	testDecodeFramework(t, "replyDouble_m0_001", -0.001)
	testDecodeFramework(t, "replyDouble_m128_0", -128.0)
	testDecodeFramework(t, "replyDouble_m129_0", -129.0)
	testDecodeFramework(t, "replyDouble_m32768_0", -32768.0)
}

func TestDoubleEncode(t *testing.T) {
	testJavaDecode(t, "argDouble_0_0", 0.0)
	testJavaDecode(t, "argDouble_0_001", 0.001)
	testJavaDecode(t, "argDouble_1_0", 1.0)
	testJavaDecode(t, "argDouble_127_0", 127.0)
	testJavaDecode(t, "argDouble_128_0", 128.0)
	testJavaDecode(t, "argDouble_2_0", 2.0)
	testJavaDecode(t, "argDouble_3_14159", 3.14159)
	testJavaDecode(t, "argDouble_32767_0", 32767.0)
	testJavaDecode(t, "argDouble_65_536", 65.536)
	testJavaDecode(t, "argDouble_m0_001", -0.001)
	testJavaDecode(t, "argDouble_m128_0", -128.0)
	testJavaDecode(t, "argDouble_m129_0", -129.0)
	testJavaDecode(t, "argDouble_m32768_0", -32768.0)
}

func TestDoublePrtEncode(t *testing.T) {
	var (
		f0 float32 = 3.14
		f1 float64 = 1024.1024102410241024
	)
	testSimpleEncode(t, &f0)
	testSimpleEncode(t, &f1)
}
