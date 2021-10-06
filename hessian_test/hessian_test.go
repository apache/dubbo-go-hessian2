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
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/stretchr/testify/assert"
)

type Case struct {
	A string
	B int
}

type CaseA struct {
	A string
	B int
	C Case
}

type CaseB struct {
	A string
	B CaseA
}

func (c *CaseB) JavaClassName() string {
	return "com.test.caseb"
}

func (c CaseA) JavaClassName() string {
	return "com.test.casea"
}

// JavaClassName  java fully qualified path
func (c Case) JavaClassName() string {
	return "com.test.case"
}

func doTestHessianEncodeHeader(t *testing.T, packageType hessian.PackageType, responseStatus byte, body interface{}) ([]byte, error) {
	hessian.RegisterPOJO(&Case{})
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

func doTestResponse(t *testing.T, packageType hessian.PackageType, responseStatus byte, body interface{}, decodedResponse *hessian.Response, assertFunc func()) {
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

	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)
	t.Log(decodedResponse)

	if assertFunc != nil {
		assertFunc()
		return
	}

	if h.ResponseStatus != hessian.Zero && h.ResponseStatus != hessian.Response_OK {
		assert.Equal(t, "java exception:"+body.(string), decodedResponse.Exception.Error())
		return
	}

	in, _ := hessian.EnsureInterface(hessian.UnpackPtrValue(hessian.EnsurePackValue(body)), nil)
	out, _ := hessian.EnsureInterface(hessian.UnpackPtrValue(hessian.EnsurePackValue(decodedResponse.RspObj)), nil)
	assert.Equal(t, in, out)
}

func TestResponse(t *testing.T) {
	caseObj := Case{A: "a", B: 1}
	decodedResponse := &hessian.Response{}

	arr := []*Case{&caseObj}
	var arrRes []interface{}
	decodedResponse.RspObj = &arrRes
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, arr, decodedResponse, func() {
		assert.Equal(t, 1, len(arrRes))
		assert.Equal(t, &caseObj, arrRes[0])
	})

	decodedResponse.RspObj = &Case{}
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, &Case{A: "a", B: 1}, decodedResponse, nil)

	s := "ok!!!!!"
	strObj := ""
	decodedResponse.RspObj = &strObj
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, s, decodedResponse, nil)

	var intObj int64
	decodedResponse.RspObj = &intObj
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, int64(3), decodedResponse, nil)

	boolObj := false
	decodedResponse.RspObj = &boolObj
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, true, decodedResponse, nil)

	strObj = ""
	decodedResponse.RspObj = &strObj
	doTestResponse(t, hessian.PackageResponse, hessian.Response_SERVER_ERROR, "error!!!!!", decodedResponse, nil)

	mapObj := map[string][]*Case{"key": {&caseObj}}
	mapRes := map[interface{}]interface{}{}
	decodedResponse.RspObj = &mapRes
	doTestResponse(t, hessian.PackageResponse, hessian.Response_OK, mapObj, decodedResponse, func() {
		c, ok := mapRes["key"]
		if !ok {
			assert.FailNow(t, "no key in decoded response map")
		}

		mapValueArr, ok := c.([]*Case)
		if !ok {
			assert.FailNow(t, "invalid decoded response map value", "expect []*Case, but get %v", reflect.TypeOf(c))
		}
		assert.Equal(t, 1, len(mapValueArr))
		assert.Equal(t, &caseObj, mapValueArr[0])
	})
}

func doTestRequest(t *testing.T, packageType hessian.PackageType, responseStatus byte, body interface{}) {
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

func TestRequest(t *testing.T) {
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a"})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a", 3})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a", true})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a", 3, true})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{3.2, true})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a", 3, true, &Case{A: "a", B: 3}})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{"a", 3, true, []*Case{{A: "a", B: 3}}})
	doTestRequest(t, hessian.PackageRequest, hessian.Zero, []interface{}{map[string][]*Case{"key": {{A: "a", B: 3}}}})
}

func TestHessianCodec_ReadAttachments(t *testing.T) {
	body := &hessian.Response{
		RspObj:      &CaseB{A: "A", B: CaseA{A: "a", B: 1, C: Case{A: "c", B: 2}}},
		Exception:   nil,
		Attachments: map[string]string{hessian.DUBBO_VERSION_KEY: "2.6.4"},
	}
	resp, err := doTestHessianEncodeHeader(t, hessian.PackageResponse, hessian.Response_OK, body)
	assert.NoError(t, err)
	hessian.UnRegisterPOJOs(&CaseB{}, &CaseA{})
	codecR1 := hessian.NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))
	codecR2 := hessian.NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))
	h := &hessian.DubboHeader{}
	assert.NoError(t, codecR1.ReadHeader(h))
	t.Log(h)
	assert.NoError(t, codecR2.ReadHeader(h))
	t.Log(h)

	err = codecR1.ReadBody(body)
	assert.Equal(t, "can not find go type name com.test.caseb in registry", err.Error())
	attrs, err := codecR2.ReadAttachments()
	assert.NoError(t, err)
	assert.Equal(t, "2.6.4", attrs[hessian.DUBBO_VERSION_KEY])

	t.Log(attrs)
}
