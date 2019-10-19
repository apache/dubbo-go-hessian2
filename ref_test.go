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
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

type circular struct {
	Num      int
	Previous *circular
	Next     *circular
}

func (circular) JavaClassName() string {
	return "circular"
}

func TestRef(t *testing.T) {
	c := &circular{}
	c.Num = 12345
	c.Previous = c
	c.Next = c

	e := NewEncoder()
	err := e.Encode(c)
	if err != nil {
		panic(err)
	}

	bytes := e.Buffer()
	t.Logf("circular bytes hex: %x, string: %s", bytes, string(bytes))
	res, err := NewDecoder(bytes).Decode()
	if err != nil {
		panic(err)
	}
	c1, ok := res.(*circular)
	if !ok {
		t.Fatalf("res:%T is not of type circular", c1)
	}
	t.Logf("encode object: %+v, decode object: %+v", c, c1)
	if c.Num != c1.Num {
		t.Errorf("encoded value %d != decoded value %d", c.Num, c1.Num)
	}
	if c1.Previous != c1.Next {
		t.Errorf("decoded value previous %p != decoded value next %p", c1.Previous, c1.Next)
	}
}

type personT struct {
	Name      string
	Relations []*personT
	Parent    *personT
	Marks     *map[string]*personT
	Tags      map[string]*personT
}

func (personT) JavaClassName() string {
	return "person"
}

func logRefObject(t *testing.T, n string, i interface{}) {
	t.Logf("ref obj[%s]: %p, %v", n, i, i)
}

func doTestRef(t *testing.T, c interface{}, name string) interface{} {
	e := NewEncoder()
	err := e.Encode(c)
	if err != nil {
		assert.FailNowf(t, "failed to encode", "error: %v", err)
	}
	bytes := e.Buffer()

	t.Logf("%s ref bytes: %s", name, string(bytes))
	t.Logf("%s ref bytes: %x", name, bytes)

	d := NewDecoder(bytes)
	decoded, err := EnsureInterface(d.Decode())
	if err != nil {
		assert.FailNowf(t, "failed to encode", "error: %v", err)
	}
	t.Logf("%s ref decoded: %v", name, decoded)
	return decoded
}

func buildComplexLevelPerson() *personT {
	p1 := &personT{Name: "p1"}
	p2 := &personT{Name: "p2"}
	p3 := &personT{Name: "p3"}
	p4 := &personT{Name: "p4"}
	p5 := &personT{Name: "p5"}
	p6 := &personT{Name: "p6"}

	p1.Parent = p2
	p2.Parent = p3
	p3.Parent = p4

	relations := []*personT{p5, p6}
	p3.Relations = relations
	p4.Relations = relations

	marks := &map[string]*personT{
		"beautiful": p1,
		"tall":      p2,
		"fat":       p3,
	}
	p4.Marks = marks
	p5.Marks = marks

	tags := map[string]*personT{
		"man":   p3,
		"woman": p4,
	}
	p5.Tags = tags
	p6.Tags = tags

	return p1
}

func TestComplexLevelRef(t *testing.T) {
	p1 := buildComplexLevelPerson()
	decoded := doTestRef(t, p1, "person")

	t.Logf("decoded object type: %v", reflect.TypeOf(decoded))
	d1, ok := decoded.(*personT)
	if !ok {
		assert.FailNow(t, "decode object is not a pointer of person")
	}
	logRefObject(t, "d1", d1)

	d2 := d1.Parent
	assert.NotNil(t, d2)
	logRefObject(t, "d2", d2)

	d3 := d2.Parent
	assert.NotNil(t, d3)
	logRefObject(t, "d3", d3)

	d4 := d3.Parent
	logRefObject(t, "d4", d4)

	assert.Equal(t, 2, len(d3.Relations))
	if len(d3.Relations) != 2 {
		assert.FailNow(t, "the length of relation array should be 2")
	}
	d5 := d3.Relations[0]
	logRefObject(t, "d5", d5)
	d6 := d3.Relations[1]
	logRefObject(t, "d6", d6)

	assert.NotNil(t, d4)
	assert.NotNil(t, d5)
	assert.NotNil(t, d6)

	assert.Equal(t, "p1", d1.Name)
	assert.Equal(t, "p2", d2.Name)
	assert.Equal(t, "p3", d3.Name)
	assert.Equal(t, "p4", d4.Name)
	assert.Equal(t, "p5", d5.Name)
	assert.Equal(t, "p6", d6.Name)

	//value equal
	assert.True(t, reflect.DeepEqual(d3.Relations, d4.Relations))

	if d4.Marks == nil {
		assert.FailNow(t, "d4.Marks should not be nil")
	}

	assert.Equal(t, 3, len(*d4.Marks))
	assert.True(t, AddrEqual(d4.Marks, d5.Marks))
	assert.True(t, AddrEqual(d1, (*d4.Marks)["beautiful"]))
	assert.True(t, AddrEqual(d2, (*d4.Marks)["tall"]))
	assert.True(t, AddrEqual(d3, (*d4.Marks)["fat"]))

	if d5.Tags == nil {
		assert.FailNow(t, "d5.Tags should not be nil")
	}
	assert.Equal(t, 2, len(d5.Tags))
	assert.True(t, reflect.DeepEqual(d5.Tags, d6.Tags))
	assert.False(t, AddrEqual(d5.Tags, d6.Tags))
	assert.True(t, AddrEqual(d3, d5.Tags["man"]))
	assert.True(t, AddrEqual(d4, d5.Tags["woman"]))
}
