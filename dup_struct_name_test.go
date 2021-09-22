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
	"bufio"
	"bytes"
	dup "dup_struct_name"
	"github.com/stretchr/testify/assert"
	"testing"
)

type CaseZ struct {
	Name string
}

func (CaseZ) JavaClassName() string {
	return "com.test.caseZ"
}

func TestDupStructNameRequest(t *testing.T) {
	RegisterPOJO(&dup.CaseZ{})
	RegisterPOJO(&CaseZ{})

	packageType := PackageRequest
	responseStatus := Zero
	var body interface{}
	body = []interface{}{"a"}
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)
	assert.Nil(t, err)

	codecR := NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &DubboHeader{}
	err = codecR.ReadHeader(h)
	assert.Nil(t, err)
	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(PackageRequest|PackageResponse|PackageHeartbeat))
	assert.Equal(t, int64(1), h.ID)
	assert.Equal(t, responseStatus, h.ResponseStatus)

	c := make([]interface{}, 7)
	err = codecR.ReadBody(c)
	assert.Nil(t, err)
	t.Log(c)
	assert.True(t, len(body.([]interface{})) == len(c[5].([]interface{})))
}

func TestDupStructNameResponse(t *testing.T) {
	RegisterPOJO(&dup.CaseZ{})
	RegisterPOJO(&CaseZ{})

	packageType := PackageResponse
	responseStatus := Response_OK
	var body interface{}
	body = &CaseZ{Name: "TestDupStructNameResponse"}
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)
	assert.Nil(t, err)

	codecR := NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &DubboHeader{}
	err = codecR.ReadHeader(h)
	assert.Nil(t, err)

	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(PackageRequest|PackageResponse|PackageHeartbeat))
	assert.Equal(t, int64(1), h.ID)
	assert.Equal(t, responseStatus, h.ResponseStatus)

	decodedResponse := &Response{}
	decodedResponse.RspObj = &CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)
	t.Log(decodedResponse)

	if h.ResponseStatus != Zero && h.ResponseStatus != Response_OK {
		assert.Equal(t, "java exception:"+body.(string), decodedResponse.Exception.Error())
		return
	}

	in, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(body)), nil)
	out, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(decodedResponse.RspObj)), nil)
	assert.Equal(t, in, out)
}
