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
	ba := &hessian.ByteArray{Values: []byte{byte(1), byte(100), byte(200)}}
	e.Encode(ba)
	return e.Buffer()
}

// nolint
func ShortArray() []byte {
	e := hessian.NewEncoder()
	sa := &hessian.ShortArray{Values: []int16{1, 100, 10000}}
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
	ba := &hessian.BooleanArray{Values: []bool{true, false, true}}
	e.Encode(ba)
	return e.Buffer()
}

// nolint
func CharacterArray() []byte {
	e := hessian.NewEncoder()
	ca := &hessian.CharacterArray{Values: "hello world"}
	e.Encode(ca)
	return e.Buffer()
}

// nolint
func FloatArray() []byte {
	e := hessian.NewEncoder()
	fa := &hessian.FloatArray{Values: []float32{1.0, 100.0, 10000.1}}
	e.Encode(fa)
	return e.Buffer()
}

// nolint
func DoubleArray() []byte {
	e := hessian.NewEncoder()
	da := &hessian.DoubleArray{Values: []float64{1.0, 100.0, 10000.1}}
	e.Encode(da)
	return e.Buffer()
}
