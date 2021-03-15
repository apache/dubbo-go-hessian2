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
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_util"
)

func TestJavaUtil(t *testing.T) {
	res, err := decodeJavaResponse(`customReplyUUID`, ``, false)
	if err != nil {
		t.Error(err)
		return
	}
	m := res.(map[interface{}]interface{})

	uuid1 := &java_util.UUID{LeastSigBits: int64(-7160773830801198154), MostSigBits: int64(459021424248441700)}

	resUuid1 := m["uuid1"]
	resUuid1String := m["uuid1_string"]
	resUuid2 := m["uuid2"]
	resUuid2String := m["uuid2_string"]

	assert.NotNil(t, resUuid1)
	assert.NotNil(t, resUuid1String)
	assert.NotNil(t, resUuid2)
	assert.NotNil(t, resUuid2String)

	assert.Equal(t, uuid1, resUuid1)
	assert.Equal(t, uuid1.String(), resUuid1String)
	assert.Equal(t, (resUuid2.(*java_util.UUID)).String(), resUuid2String)
}
