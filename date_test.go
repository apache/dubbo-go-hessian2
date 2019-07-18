// Copyright 2016-2019 Alex Stocks, Xinge Gao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hessian

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	RegisterPOJO(&DateDemo{})
}

type DateDemo struct {
	Name string
	Date time.Time
}

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
	r, e := decodeJavaResponse(method, "")
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
		e   *Encoder
		d   *Decoder
		res interface{}
	)
	date := DateDemo{
		Name: "s",
		Date: ZeroDate,
	}
	e = NewEncoder()
	e.Encode(date)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}
	d = NewDecoder(e.Buffer())
	res, _ = d.Decode()
	assert.Equal(t, ZeroDate, res.(*DateDemo).Date)
}

func TestDateNulJavaDecode(t *testing.T) {
	date := DateDemo{
		Name: "s",
		Date: ZeroDate,
	}
	testJavaDecode(t, "customArgTypedFixedList_DateNull", date)
}

func TestDateNullDecode(t *testing.T) {

	doTestDateNull(t, "customReplyTypedFixedDateNull")
}

func doTestDateNull(t *testing.T, method string) {
	testDecodeFrameworkFunc(t, method, func(r interface{}) {
		t.Logf("%#v", r)
		assert.Equal(t, ZeroDate, r.(*DateDemo).Date)
	})
}
