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

package hessian_test

import (
	"bufio"
	"bytes"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-hessian2"
	dupclass "github.com/apache/dubbo-go-hessian2/hessian_test/hessian_test"
)

import (
	"github.com/stretchr/testify/assert"
)

const (
	ExpectedErrorMsg = "reflect.Set: value of type hessian_test.CaseZ is not assignable to type hessian_test.CaseZ"
)

type CaseZ struct {
	Name string
}

func (CaseZ) JavaClassName() string {
	return "com.test.caseZ"
}

func doTestHessianEncodeHeader(t *testing.T, packageType hessian.PackageType, responseStatus byte, body interface{}) ([]byte, error) {
	codecW := hessian.NewHessianCodec(nil)
	resp, err := codecW.Write(hessian.Service{
		Path:      "test",
		Interface: "ITest",
		Version:   "v1.0",
		Method:    "test",
		Timeout:   time.Second * 10,
	}, hessian.DubboHeader{
		SerialID:       2,
		Type:           packageType,
		ID:             1,
		ResponseStatus: responseStatus,
	}, body)
	assert.Nil(t, err)
	return resp, err
}

func TestDupStructNameRequest(t *testing.T) {
	hessian.RegisterPOJO(&dupclass.CaseZ{})
	hessian.RegisterPOJO(&CaseZ{})

	packageType := hessian.PackageRequest
	responseStatus := hessian.Zero
	var body interface{}
	body = []interface{}{"a"}
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)
	assert.Nil(t, err)

	codecR := hessian.NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &hessian.DubboHeader{}
	err = codecR.ReadHeader(h)
	assert.Nil(t, err)
	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(hessian.PackageRequest|hessian.PackageResponse|hessian.PackageHeartbeat))
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
				assert.Equal(t, ExpectedErrorMsg, errStr)
			}
		}
	}()

	var body interface{}
	body = &CaseZ{Name: "TestDupStructNameResponse"}
	err, codecR, h := doTestHeader(t, body)
	assert.Nil(t, err)

	decodedResponse := &hessian.Response{}
	decodedResponse.RspObj = &dupclass.CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.NotNil(t, err)
	assert.Equal(t, ExpectedErrorMsg, err.Error())

	decodedResponse = &hessian.Response{}
	decodedResponse.RspObj = &CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)

	checkResponseBody(t, decodedResponse, h, body)
}

func TestDupStructNameResponse2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if errStr, ok := err.(string); ok {
				assert.Equal(t, ExpectedErrorMsg, errStr)
			}
		}
	}()

	var body interface{}
	body = &dupclass.CaseZ{Name: "TestDupStructNameResponse"}
	err, codecR, h := doTestHeader(t, body)
	assert.Nil(t, err)

	decodedResponse := &hessian.Response{}
	decodedResponse.RspObj = &CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.NotNil(t, err)
	assert.Equal(t, ExpectedErrorMsg, err.Error())

	decodedResponse = &hessian.Response{}
	decodedResponse.RspObj = &dupclass.CaseZ{}
	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)

	checkResponseBody(t, decodedResponse, h, body)
}

func doTestHeader(t *testing.T, body interface{}) (error, *hessian.HessianCodec, *hessian.DubboHeader) {
	hessian.RegisterPOJO(&dupclass.CaseZ{})
	hessian.RegisterPOJO(&CaseZ{})

	packageType := hessian.PackageResponse
	responseStatus := hessian.Response_OK
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)
	assert.Nil(t, err)

	codecR := hessian.NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &hessian.DubboHeader{}
	err = codecR.ReadHeader(h)
	assert.Nil(t, err)

	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(hessian.PackageRequest|hessian.PackageResponse|hessian.PackageHeartbeat))
	assert.Equal(t, int64(1), h.ID)
	assert.Equal(t, responseStatus, h.ResponseStatus)
	return err, codecR, h
}

func checkResponseBody(t *testing.T, decodedResponse *hessian.Response, h *hessian.DubboHeader, body interface{}) {
	t.Log(decodedResponse)

	if h.ResponseStatus != hessian.Zero && h.ResponseStatus != hessian.Response_OK {
		assert.Equal(t, "java exception:"+body.(string), decodedResponse.Exception.Error())
		return
	}

	in, _ := hessian.EnsureInterface(hessian.UnpackPtrValue(hessian.EnsurePackValue(body)), nil)
	out, _ := hessian.EnsureInterface(hessian.UnpackPtrValue(hessian.EnsurePackValue(decodedResponse.RspObj)), nil)
	assert.Equal(t, in, out)
}

func TestDuplicatedClassGetGoType(t *testing.T) {
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/hessian_test_test/hessian_test.CaseZ", hessian.GetGoType(&CaseZ{}))
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/hessian_test/hessian_test/hessian_test.CaseZ", hessian.GetGoType(&dupclass.CaseZ{}))
}
