// Copyright (c) 2016 ~ 2019, Alex Stocks.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Unit test for decoding hessian2 based on official api. One can find api
// doc on http://javadoc4.caucho.com/com/caucho/hessian/test/TestHessian2.html
package hessian

import (
	"net/http"
	"bytes"
	"log"
	"io/ioutil"
	"encoding/binary"
	"testing"
	"reflect"
)

func encodeCall(method string) []byte {
	b := []byte{'c', 2, 0, 'm'}
	bMethodLength := make([]byte, 2)
	binary.BigEndian.PutUint16(bMethodLength, uint16(len(method)))
	b = append(b, bMethodLength...)
	b = append(b, method...)
	b = append(b, 'z')
	return b
}

func sendRequest(method string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest(
		"POST",
		"http://hessian.caucho.com/test/test",
		bytes.NewReader(encodeCall(method)),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body[4:]  // skip H 0x02 0x00
}

func decodeResponse(method string) (interface{}, error) {
	b := sendRequest(method)
	d := NewDecoder(b)
	r, e := d.Decode()
	if e != nil {
		return nil, e
	}
	return r, nil
}

func testDoubleFramework(t *testing.T, method string, expected float64) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(float64)
	if ! ok {
		t.Errorf("%s: %v is not double", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestDouble(t *testing.T) {
	testDoubleFramework(t, "replyDouble_0_0", 0.0)
	testDoubleFramework(t, "replyDouble_0_001", 0.001)
	testDoubleFramework(t, "replyDouble_1_0", 1.0)
	testDoubleFramework(t, "replyDouble_127_0", 127.0)
	testDoubleFramework(t, "replyDouble_128_0", 128.0)
	testDoubleFramework(t, "replyDouble_2_0", 2.0)
	testDoubleFramework(t, "replyDouble_3_14159", 3.14159)
	testDoubleFramework(t, "replyDouble_32767_0", 32767.0)
	testDoubleFramework(t, "replyDouble_65_536", 65.536)
	testDoubleFramework(t, "replyDouble_m0_001", -0.001)
	testDoubleFramework(t, "replyDouble_m128_0", -128.0)
	testDoubleFramework(t, "replyDouble_m129_0", -129.0)
	testDoubleFramework(t, "replyDouble_m32768_0", -32768.0)
}

func testBooleanFramework(t *testing.T, method string, expected bool) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(bool)
	if ! ok {
		t.Errorf("%s: %v is not bool", method, r)
		return
	}

	if ok && v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestBoolean(t *testing.T) {
	testBooleanFramework(t, "replyFalse", false)
	testBooleanFramework(t, "replyTrue", true)
}

func testIntFramework(t *testing.T, method string, expected int32) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(int32)
	if ! ok {
		t.Errorf("%s: %v is not int", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestInt(t *testing.T) {
	testIntFramework(t, "replyInt_0", 0)
	testIntFramework(t, "replyInt_0x30", 0x30)
	testIntFramework(t, "replyInt_0x3ffff", 0x3ffff)
	testIntFramework(t, "replyInt_0x40000", 0x40000)
	testIntFramework(t, "replyInt_0x7ff", 0x7ff)
	testIntFramework(t, "replyInt_0x7fffffff", 0x7fffffff)
	testIntFramework(t, "replyInt_0x800", 0x800)
	testIntFramework(t, "replyInt_1", 1)
	testIntFramework(t, "replyInt_47", 47)
	testIntFramework(t, "replyInt_m0x40000", -0x40000)
	testIntFramework(t, "replyInt_m0x40001", -0x40001)
	testIntFramework(t, "replyInt_m0x800", -0x800)
	testIntFramework(t, "replyInt_m0x80000000", -0x80000000)
	testIntFramework(t, "replyInt_m0x801", -0x801)
	testIntFramework(t, "replyInt_m16", -16)
	testIntFramework(t, "replyInt_m17", -17)
}

func testLongFramework(t *testing.T, method string, expected int64) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	v, ok := r.(int64)
	if ! ok {
		t.Errorf("%s: %v is not long", method, r)
		return
	}

	if v != expected {
		t.Errorf("%s: got %v, wanted %v", method, v, expected)
	}
}

func TestLong(t *testing.T) {
	testLongFramework(t, "replyLong_0", 0)
	testLongFramework(t, "replyLong_0x10", 0x10)
	testLongFramework(t, "replyLong_0x3ffff", 0x3ffff)
	testLongFramework(t, "replyLong_0x40000", 0x40000)
	testLongFramework(t, "replyLong_0x7ff", 0x7ff)
	testLongFramework(t, "replyLong_0x7fffffff", 0x7fffffff)
	testLongFramework(t, "replyLong_0x800", 0x800)
	testLongFramework(t, "replyLong_1", 1)
	testLongFramework(t, "replyLong_15", 15)
	testLongFramework(t, "replyLong_m0x40000", -0x40000)
	testLongFramework(t, "replyLong_m0x40001", -0x40001)
	testLongFramework(t, "replyLong_m0x800", -0x800)
	testLongFramework(t, "replyLong_m0x80000000", -0x80000000)
	testLongFramework(t, "replyLong_m0x80000001", -0x80000001)
	testLongFramework(t, "replyLong_m0x801", -0x801)
	testLongFramework(t, "replyLong_m8", -8)
	testLongFramework(t, "replyLong_m9", -9)
}

func TestNull(t *testing.T) {
	r, e := decodeResponse("replyNull")
	if e != nil {
		t.Errorf("replyNull: decode fail with error %v", e)
		return
	}

	if reflect.TypeOf(r) != nil {  // detect nil interface, not only nil value
		t.Errorf("replyNull: %v is not null", r)
	}
}
