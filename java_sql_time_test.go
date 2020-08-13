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
	"time"
)

func init() {
	SetJavaSqlTimeSerialize(&Date{})
	SetJavaSqlTimeSerialize(&Time{})
}
func TestJavaSqlTimeEncode(t *testing.T) {
	Bdate := "1997-01-01 13:15:46"
	location, _ := time.ParseInLocation("2006-01-02 15:04:05", Bdate, time.Local)
	time1 := Time{Time: location}
	t.Log(time1.UnixNano() / 1000000)
	testJavaDecode(t, "javaSql_encode_time", time1)
}

func TestJavaSqlTimeDecode(t *testing.T) {
	Bdate := "1997-01-01 13:15:46"
	location, _ := time.ParseInLocation("2006-01-02 15:04:05", Bdate, time.Local)
	time1 := Time{Time: location}
	//t.Log(time1.UnixNano() / 1000000)
	//testJavaDecode(t, "javaSql_encode_time", time1)
	testDecodeFramework(t, "javaSql_decode_time", &time1)
}

func TestName(t *testing.T) {
	Bdate := "1997-01-01 13:15:46"
	location, _ := time.ParseInLocation("2006-01-02 15:04:05", Bdate, time.Local)
	time1 := Time{Time: location}
	time2 := Time{Time: location}
	r := time1
	expected := time2
	t.Log(reflect.DeepEqual(r, expected))
}
