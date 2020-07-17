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
	"github.com/apache/dubbo-go-hessian2/java8_time"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJava8Time(t *testing.T) {

	//use go decode java data(data from java encode byte[])
	doTestTime(t, "java8_Year", &java8_time.Year{Year: 2020})
	//use go and java encode, compare the encode results(string) of java and go
	doTestJava8Time(t, "java8_Year", java8_time.Year{Year: 2020})
	doTestTime(t, "java8_LocalDate", &java8_time.LocalDate{Year: 2020, Month: 6, Day: 6})
	doTestTime(t, "java8_LocalTime", &java8_time.LocalTime{Hour: 6, Minute: 6})
	doTestTime(t, "java8_LocalDateTime", &java8_time.LocalDateTime{Date: java8_time.LocalDate{Year: 2020, Month: 6, Day: 6}, Time: java8_time.LocalTime{Hour: 6, Minute: 6}})
	doTestTime(t, "java8_MonthDay", &java8_time.MonthDay{Month: 6, Day: 6})
	doTestTime(t, "java8_Duration", &java8_time.Duration{Second: 0, Nano: 0})
	doTestTime(t, "java8_Instant", &java8_time.Instant{Seconds: 100, Nanos: 0})
}

func doTestJava8Time(t *testing.T, method string, pojo POJO) {
	e := NewEncoder()
	err := e.Encode(pojo)
	// go encode
	goStr := e.buffer
	if err != nil {
		panic(err)
	}
	//java encode
	javaStr := getJavaReply(method, "")
	assert.Equal(t, goStr, javaStr)

}

func doTestTime(t *testing.T, method string, expected interface{}) {
	testDecodeFramework(t, method, expected)
}
