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

func TestPackRequest(t *testing.T) {
	bytes, err := packRequest(Service{
		Path:      "test",
		Interface: "ITest",
		Version:   "v1.0",
		Method:    "test",
		Timeout:   time.Second * 10,
	}, DubboHeader{
		SerialID: 0,
		Type:     PackageRequest,
		ID:       123,
	}, []interface{}{1, 2})

	assert.Nil(t, err)

	if bytes != nil {
		t.Logf("pack request: %s", string(bytes))
	}
}

func TestGetArgsTypeList(t *testing.T) {
	type Test struct{}
	str, err := getArgsTypeList([]interface{}{nil, 1, []int{2}, true, []bool{false}, "a", []string{"b"}, Test{}, &Test{}, []Test{}, map[string]Test{}})
	assert.NoError(t, err)
	assert.Equal(t, "VJ[JZ[ZLjava/lang/String;[Ljava/lang/String;Ljava/lang/Object;Ljava/lang/Object;[Ljava/lang/Object;Ljava/util/Map;", str)
}

func TestDescRegex(t *testing.T) {
	results := DescRegex.FindAllString("Ljava/lang/String;", -1)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "Ljava/lang/String;", results[0])

	results = DescRegex.FindAllString("Ljava/lang/String;I", -1)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "Ljava/lang/String;", results[0])
	assert.Equal(t, "I", results[1])

	results = DescRegex.FindAllString("ILjava/lang/String;", -1)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "I", results[0])
	assert.Equal(t, "Ljava/lang/String;", results[1])

	results = DescRegex.FindAllString("ILjava/lang/String;IZ", -1)
	assert.Equal(t, 4, len(results))
	assert.Equal(t, "I", results[0])
	assert.Equal(t, "Ljava/lang/String;", results[1])
	assert.Equal(t, "I", results[2])
	assert.Equal(t, "Z", results[3])

	results = DescRegex.FindAllString("[Ljava/lang/String;[I", -1)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "[Ljava/lang/String;", results[0])
	assert.Equal(t, "[I", results[1])
}
