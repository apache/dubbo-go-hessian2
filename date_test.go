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
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func init() {
	RegisterPOJO(&DateDemo{})
}

type DateDemo struct {
	Name    string
	Date    time.Time
	Dates   []**time.Time
	NilDate *time.Time
	Date1   *time.Time
	Date2   **time.Time
	Date3   ***time.Time
}

// JavaClassName  java fully qualified path
func (DateDemo) JavaClassName() string {
	return "test.model.DateDemo"
}

func TestEncDate(t *testing.T) {
	var (
		v   string
		tz  time.Time
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	v = "2014-02-09 06:15:23"
	tz, _ = time.Parse("2006-01-02 15:04:05", v)
	e.Encode(tz)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	t.Logf("decode(%s, %s) = %v, %v\n", v, tz.Local(), res, err)
}

func testDateFramework(t *testing.T, method string, expected time.Time) {
	r, e := decodeJavaResponse(method, "", false)
	if e != nil {
		t.Errorf("%s: decode fail with error %+v", method, e)
		return
	}

	v, ok := r.(time.Time)
	if !ok {
		t.Errorf("%s: %v is not date", method, r)
		return
	}

	if !v.Equal(expected) {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestDate(t *testing.T) {
	testDateFramework(t, "replyDate_0", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	testDateFramework(t, "replyDate_1", time.Date(1998, 5, 8, 9, 51, 31, 0, time.UTC))
	testDateFramework(t, "replyDate_2", time.Date(1998, 5, 8, 9, 51, 0, 0, time.UTC))
}

func TestDateEncode(t *testing.T) {
	testJavaDecode(t, "argDate_0", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	testJavaDecode(t, "argDate_1", time.Date(1998, 5, 8, 9, 51, 31, 0, time.UTC))
	testJavaDecode(t, "argDate_2", time.Date(1998, 5, 8, 9, 51, 0, 0, time.UTC))
}

func TestEncDateNull(t *testing.T) {
	var (
		v   string
		tz  time.Time
		e   *Encoder
		d   *Decoder
		res interface{}
		err error
	)
	v = "2014-02-09 06:15:23 +0800 CST"
	tz, _ = time.Parse("2006-01-02 15:04:05 +0800 CST", v)
	d1 := &tz
	d2 := &d1
	d3 := &d2

	date := DateDemo{
		Name:    "zs",
		Date:    ZeroDate,
		Dates:   []**time.Time{d2, d2},
		NilDate: nil,
		Date1:   nil,
		Date2:   d2,
		Date3:   d3,
	}
	e = NewEncoder()
	err = e.Encode(date)
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Buffer()) == 0 {
		t.Fail()
	}
	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	assert.Nil(t, err)
	assert.Equal(t, ZeroDate, res.(*DateDemo).Date)
	assert.Equal(t, 2, len(res.(*DateDemo).Dates))
	assert.Equal(t, tz.Local().String(), (*res.(*DateDemo).Dates[0]).String())
	assert.Equal(t, &ZeroDate, res.(*DateDemo).NilDate)
	assert.Equal(t, ZeroDate, *res.(*DateDemo).Date1)
	assert.Equal(t, tz.Local().String(), (*res.(*DateDemo).Date2).String())
	assert.Equal(t, tz.Local().String(), (*(*res.(*DateDemo).Date3)).String())
}

func TestDateNulJavaDecode(t *testing.T) {
	date := DateDemo{
		Name: "zs",
		Date: ZeroDate,
	}
	testJavaDecode(t, "customArgTypedFixed_DateNull", date)
}

func TestDateNilDecode(t *testing.T) {
	doTestDateNull(t, "customReplyTypedFixedDateNull")
}

func doTestDateNull(t *testing.T, method string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, ZeroDate, r.(*DateDemo).Date)
		assert.Equal(t, &ZeroDate, r.(*DateDemo).Date1)
	})
}
