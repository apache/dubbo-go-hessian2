// Copyright 2016-2019 Alex Stocks
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

package hessian

import (
	"testing"
)

func TestEncList(t *testing.T) {
	var (
		list []interface{}
		err  error
		e    *Encoder
		d    *Decoder
		res  interface{}
	)

	e = NewEncoder()
	list = []interface{}{100, 10.001, "hello", []byte{0, 2, 4, 6, 8, 10}, true, nil, false}
	e.Encode(list)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", list, res, err)
}

func TestList(t *testing.T) {
	testDecodeFramework(t, "replyTypedFixedList_0", []interface{}{})
	testDecodeFramework(t, "replyTypedFixedList_1", []interface{}{"1"})
	testDecodeFramework(t, "replyTypedFixedList_7", []interface{}{"1", "2", "3", "4", "5", "6", "7"})
	testDecodeFramework(t, "replyTypedFixedList_8", []interface{}{"1", "2", "3", "4", "5", "6", "7", "8"})
	testDecodeFramework(t, "replyUntypedFixedList_0", []interface{}{})
	testDecodeFramework(t, "replyUntypedFixedList_1", []interface{}{"1"})
	testDecodeFramework(t, "replyUntypedFixedList_7", []interface{}{"1", "2", "3", "4", "5", "6", "7"})
	testDecodeFramework(t, "replyUntypedFixedList_8", []interface{}{"1", "2", "3", "4", "5", "6", "7", "8"})
}
