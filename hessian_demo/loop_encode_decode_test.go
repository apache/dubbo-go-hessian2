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

package hessian_demo

import (
	"reflect"
	"testing"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name string
	Age  int
}

func (User) JavaClassName() string {
	return "hessian.demo.User"
}

type Location struct {
	Address  string
	Postcode string
}

func (Location) JavaClassName() string {
	return "hessian.demo.Location"
}

func TestLoopDecode(t *testing.T) {
	u := &User{
		Name: "wongoo",
		Age:  18,
	}
	loc := &Location{
		Address:  "xiamen",
		Postcode: "361000",
	}

	// demo a bytes buffer from client to server
	var bytes []byte
	e := hessian.NewEncoder()

	// -------- encode 1
	if err := e.Encode(u); err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, e.Buffer()...)

	// -------- encode 2
	e.Clean()
	if err := e.Encode(12345); err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, e.Buffer()...)

	// -------- encode 3
	e.Clean()
	if err := e.Encode(loc); err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, e.Buffer()...)

	// -------- encode 4
	e.Clean()
	if err := e.Encode("hello"); err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, e.Buffer()...)

	// demo a decoder to decode buffer from client
	d := hessian.NewDecoder(bytes)

	// -------- decode 1
	res, err := d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, reflect.DeepEqual(u, res))

	// -------- decode 2
	d.Clean()
	res, err = d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int64(12345), res)

	// -------- decode 3
	d.Clean()
	res, err = d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, reflect.DeepEqual(loc, res))

	// -------- decode 4
	d.Clean()
	res, err = d.Decode()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", res)
}
