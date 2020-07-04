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

	doTestJava8Time(t, "java8_Year", java8_time.Year{Year: 2020})
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
