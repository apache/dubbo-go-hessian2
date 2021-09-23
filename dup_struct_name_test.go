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
	"testing"
)

import (
	dup "dubbo-go-hessian2-dup-struct-name-test"
)

import (
	"github.com/stretchr/testify/assert"
)

const (
	EXPECTED_ERROR_MSG = "reflect.Set: value of type hessian.CaseZ is not assignable to type hessian.CaseZ"
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
	defer func() {
		if err := recover(); err != nil {
			if errStr, ok := err.(string); ok {
				assert.Equal(t, EXPECTED_ERROR_MSG, errStr)
			}
		}
	}()

	var body interface{}
	body = &CaseZ{Name: "TestDupStructNameResponse"}
	err, codecR, h := doTestHeader(t, body)

	decodedResponse := &Response{}
	decodedResponse.RspObj = &dup.CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.NotNil(t, err)
	assert.Equal(t, EXPECTED_ERROR_MSG, err.Error())

	decodedResponse = &Response{}
	decodedResponse.RspObj = &CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)

	checkResponseBody(t, decodedResponse, h, body)
}

func TestDupStructNameResponse2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if errStr, ok := err.(string); ok {
				assert.Equal(t, EXPECTED_ERROR_MSG, errStr)
			}
		}
	}()

	var body interface{}
	body = &dup.CaseZ{Name: "TestDupStructNameResponse"}
	err, codecR, h := doTestHeader(t, body)

	decodedResponse := &Response{}
	decodedResponse.RspObj = &CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.NotNil(t, err)
	assert.Equal(t, EXPECTED_ERROR_MSG, err.Error())

	decodedResponse = &Response{}
	decodedResponse.RspObj = &dup.CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)

	checkResponseBody(t, decodedResponse, h, body)
}

func doTestHeader(t *testing.T, body interface{}) (error, *HessianCodec, *DubboHeader) {
	RegisterPOJO(&dup.CaseZ{})
	RegisterPOJO(&CaseZ{})

	packageType := PackageResponse
	responseStatus := Response_OK
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
	return err, codecR, h
}

func checkResponseBody(t *testing.T, decodedResponse *Response, h *DubboHeader, body interface{}) {
	t.Log(decodedResponse)

	if h.ResponseStatus != Zero && h.ResponseStatus != Response_OK {
		assert.Equal(t, "java exception:"+body.(string), decodedResponse.Exception.Error())
		return
	}

	in, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(body)), nil)
	out, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(decodedResponse.RspObj)), nil)
	assert.Equal(t, in, out)
}
