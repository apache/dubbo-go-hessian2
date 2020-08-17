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

import (
	"github.com/apache/dubbo-go-hessian2/java_sql_time"
)

func init() {
	SetJavaSqlTimeSerialize(&java_sql_time.Date{})
	SetJavaSqlTimeSerialize(&java_sql_time.Time{})
}

// test local time between go and java
// go encode
// java decode
func TestJavaSqlTimeEncode(t *testing.T) {
	sqlTime := time.Date(1997, 1, 1, 13, 15, 46, 0, time.UTC)
	testSqlTime := java_sql_time.Time{Time: sqlTime}
	testJavaDecode(t, "javaSql_encode_time", testSqlTime)

	sqlDate := time.Date(2020, 8, 9, 0, 0, 0, 0, time.UTC)
	testSqlDate := java_sql_time.Date{Time: sqlDate}
	testJavaDecode(t, "javaSql_encode_date", &testSqlDate)
}

// test local time between go and java
// java encode
// go decode
func TestJavaSqlTimeDecode(t *testing.T) {
	sqlTime := time.Date(1997, 1, 1, 13, 15, 46, 0, time.UTC)
	testSqlTime := java_sql_time.Time{Time: sqlTime}
	testDecodeJavaSqlTime(t, "javaSql_decode_time", &testSqlTime)

	sqlDate := time.Date(2020, 8, 9, 0, 0, 0, 0, time.UTC)
	testDateTime := java_sql_time.Date{Time: sqlDate}
	testDecodeJavaSqlTime(t, "javaSql_decode_date", &testDateTime)
}

func testDecodeJavaSqlTime(t *testing.T, method string, expected java_sql_time.JavaSqlTime) {
	r, e := decodeJavaResponse(method, "", false)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}
	resultSqlTime, ok := r.(java_sql_time.JavaSqlTime)
	if !ok {
		t.Errorf("got error type:%v", r)
	}
	assert.Equal(t, resultSqlTime.GetTime().UnixNano(), expected.GetTime().UnixNano())
}

// test local time between go and go
// go encode
// go decode
func TestJavaSqlTimeWithGo(t *testing.T) {
	location, _ := time.ParseInLocation("2006-01-02 15:04:05", "1997-01-01 13:15:46", time.Local)
	sqlTime := java_sql_time.Time{Time: location}
	e := NewEncoder()
	e.Encode(&sqlTime)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}
	d := NewDecoder(e.Buffer())
	resultSqlTime, _ := d.Decode()
	assert.Equal(t, &sqlTime, resultSqlTime)

	location, _ = time.ParseInLocation("2006-01-02 15:04:05", "2020-08-09 00:00:00", time.Local)
	sqlDate := java_sql_time.Date{Time: location}
	e = NewEncoder()
	e.Encode(&sqlDate)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}
	d = NewDecoder(e.Buffer())
	resultSqlDate, _ := d.Decode()
	assert.Equal(t, &sqlDate, resultSqlDate)
}
