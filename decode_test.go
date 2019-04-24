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
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
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
	return body[4:] // skip H 0x02 0x00
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
