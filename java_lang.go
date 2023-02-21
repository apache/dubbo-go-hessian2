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

import "reflect"

func init() {
	var a int32 = 1
	registerPOJOTypeMapping("java.lang.Integer", "*int32", reflect.TypeOf(&a), &a)

	var b int64 = 1
	registerPOJOTypeMapping("java.lang.Long", "*int64", reflect.TypeOf(&b), &b)

	var c = true
	registerPOJOTypeMapping("java.lang.Boolean", "*bool", reflect.TypeOf(&c), &c)

	var d int16 = 1
	registerPOJOTypeMapping("java.lang.Short", "*int16", reflect.TypeOf(&d), &d)

	var e byte = 'a'
	registerPOJOTypeMapping("java.lang.Byte", "*uint8", reflect.TypeOf(&e), &e)

	var f float32 = 1.0
	registerPOJOTypeMapping("java.lang.Float", "*float32", reflect.TypeOf(&f), &f)

	var g = 1.0
	registerPOJOTypeMapping("java.lang.Double", "*float64", reflect.TypeOf(&g), &g)

	var h = 'a'
	registerPOJOTypeMapping("java.lang.Character", "*hessian.Rune", reflect.TypeOf(&h), &h)
}

// Rune is an alias for rune, so that to get the correct runtime type of rune.
// The runtime type of rune is int32, which is not expected.
type Rune rune

var (
	_typeOfRune = reflect.TypeOf(Rune(0))
)
