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

package testfuncs

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

// nolint
func ByteArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c byte = 1, 100, 200
	ba := []*byte{&a, &b, &c}
	e.Encode(ba)
	return e.Buffer()
}

// nolint
func ShortArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c int16 = 1, 100, 10000
	sa := []*int16{&a, &b, &c}
	e.Encode(sa)
	return e.Buffer()
}

// nolint
func IntegerArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c int32 = 1, 100, 10000
	ia := []*int32{&a, &b, &c}
	e.Encode(ia)
	return e.Buffer()
}

// nolint
func LongArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c int64 = 1, 100, 10000
	la := []*int64{&a, &b, &c}
	e.Encode(la)
	return e.Buffer()
}

// nolint
func BooleanArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c = true, false, true
	ba := []*bool{&a, &b, &c}
	e.Encode(ba)
	return e.Buffer()
}

// nolint
func CharacterArray() []byte {
	e := hessian.NewEncoder()
	var r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11 hessian.Rune = 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'
	ca := []*hessian.Rune{&r1, &r2, &r3, &r4, &r5, &r6, &r7, &r8, &r9, &r10, &r11}
	e.Encode(ca)
	return e.Buffer()
}

// nolint
func FloatArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c float32 = 1.0, 100.0, 10000.1
	fa := []*float32{&a, &b, &c}
	e.Encode(fa)
	return e.Buffer()
}

// nolint
func DoubleArray() []byte {
	e := hessian.NewEncoder()
	var a, b, c = 1.0, 100.0, 10000.1
	da := []*float64{&a, &b, &c}
	e.Encode(da)
	return e.Buffer()
}

// nolint
func MultipleLevelA0Array() []byte {
	hessian.RegisterPOJO(&A0{})

	e := hessian.NewEncoder()
	da := [][][]*A0{
		{{new(A0), new(A0), new(A0)}, {new(A0), new(A0), new(A0), nil}},
		{{new(A0)}, {new(A0)}},
	}
	e.Encode(da)
	return e.Buffer()
}