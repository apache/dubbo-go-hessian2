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

package issue356

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

type UserInfo struct {
	Name    string
	Address map[string]string
	Family  map[string]int
}

func (UserInfo) JavaClassName() string {
	return "com.test.UserInfo"
}

func TestIssue356Case(t *testing.T) {
	info := &UserInfo{
		Name:    "test",
		Address: nil,
		Family:  nil,
	}

	hessian.RegisterPOJO(info)

	encoder := hessian.NewEncoder()
	err := encoder.Encode(info)
	if err != nil {
		t.Error(err)
		return
	}

	enBuf := encoder.Buffer()

	decoder := hessian.NewDecoder(enBuf)
	dec, err := decoder.Decode()
	assert.Nil(t, err)

	t.Log(dec)

	assert.True(t, reflect.DeepEqual(info, dec))
}
