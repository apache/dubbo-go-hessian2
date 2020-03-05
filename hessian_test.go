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
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

type Case struct {
	A string
	B int
}

//JavaClassName  java fully qualified path
func (c *Case) JavaClassName() string {
	return "com.test.case"
}

func doTestHessianEncodeHeader(t *testing.T, packageType PackageType, responseStatus byte, body interface{}) ([]byte, error) {
	RegisterPOJO(&Case{})
	codecW := NewHessianCodec(nil)
	resp, err := codecW.Write(Service{
		Path:      "test",
		Interface: "ITest",
		Version:   "v1.0",
		Method:    "test",
		Timeout:   time.Second * 10,
	}, DubboHeader{
		SerialID:       2,
		Type:           packageType,
		ID:             1,
		ResponseStatus: responseStatus,
	}, body)
	assert.Nil(t, err)
	return resp, err
}

func doTestResponse(t *testing.T, packageType PackageType, responseStatus byte, body interface{}, decodedResponse *Response, assertFunc func()) {
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)

	codecR := NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &DubboHeader{}
	err = codecR.ReadHeader(h)
	assert.Nil(t, err)

	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(PackageRequest|PackageResponse|PackageHeartbeat))
	assert.Equal(t, int64(1), h.ID)
	assert.Equal(t, responseStatus, h.ResponseStatus)

	err = codecR.ReadBody(decodedResponse)
	assert.Nil(t, err)
	t.Log(decodedResponse)

	if assertFunc != nil {
		assertFunc()
		return
	}

	if h.ResponseStatus != Zero && h.ResponseStatus != Response_OK {
		assert.Equal(t, "java exception:"+body.(string), decodedResponse.Exception.Error())
		return
	}

	in, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(body)), nil)
	out, _ := EnsureInterface(UnpackPtrValue(EnsurePackValue(decodedResponse.RspObj)), nil)
	assert.Equal(t, in, out)
}

func TestResponse(t *testing.T) {
	caseObj := Case{A: "a", B: 1}
	decodedResponse := &Response{}

	arr := []*Case{&caseObj}
	var arrRes []interface{}
	decodedResponse.RspObj = &arrRes
	doTestResponse(t, PackageResponse, Response_OK, arr, decodedResponse, func() {
		assert.Equal(t, 1, len(arrRes))
		assert.Equal(t, &caseObj, arrRes[0])
	})

	decodedResponse.RspObj = &Case{}
	doTestResponse(t, PackageResponse, Response_OK, &Case{A: "a", B: 1}, decodedResponse, nil)

	s := "ok!!!!!"
	strObj := ""
	decodedResponse.RspObj = &strObj
	doTestResponse(t, PackageResponse, Response_OK, s, decodedResponse, nil)

	var intObj int64
	decodedResponse.RspObj = &intObj
	doTestResponse(t, PackageResponse, Response_OK, int64(3), decodedResponse, nil)

	boolObj := false
	decodedResponse.RspObj = &boolObj
	doTestResponse(t, PackageResponse, Response_OK, true, decodedResponse, nil)

	strObj = ""
	decodedResponse.RspObj = &strObj
	doTestResponse(t, PackageResponse, Response_SERVER_ERROR, "error!!!!!", decodedResponse, nil)

	mapObj := map[string][]*Case{"key": {&caseObj}}
	mapRes := map[interface{}]interface{}{}
	decodedResponse.RspObj = &mapRes
	doTestResponse(t, PackageResponse, Response_OK, mapObj, decodedResponse, func() {
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

func doTestRequest(t *testing.T, packageType PackageType, responseStatus byte, body interface{}) {
	resp, err := doTestHessianEncodeHeader(t, packageType, responseStatus, body)

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

func TestRequest(t *testing.T) {
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a"})
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a", 3})
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a", true})
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a", 3, true})
	doTestRequest(t, PackageRequest, Zero, []interface{}{3.2, true})
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a", 3, true, &Case{A: "a", B: 3}})
	doTestRequest(t, PackageRequest, Zero, []interface{}{"a", 3, true, []*Case{{A: "a", B: 3}}})
	doTestRequest(t, PackageRequest, Zero, []interface{}{map[string][]*Case{"key": {{A: "a", B: 3}}}})
}
