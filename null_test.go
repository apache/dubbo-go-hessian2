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
	"reflect"
	"testing"
)

func TestEncNull(t *testing.T) {
	var e = NewEncoder()
	e.Encode(nil)
	if e.Buffer() == nil {
		t.Fail()
	}
	t.Logf("nil enc result:%s\n", string(e.buffer))
}
func testNullFramework(t *testing.T, method string) {
	r, e := decodeResponse(method)
	if e != nil {
		t.Errorf("%s: decode fail with error %+v", method, e)
		return
	}

	if reflect.TypeOf(r) != nil { // detect nil interface, not only nil value
		t.Errorf("%s: %v is not null", method, r)
	}
}

func TestNull(t *testing.T) {
	testNullFramework(t, "replyBinary_null")
	testNullFramework(t, "replyNull")
	testNullFramework(t, "replyString_null")
}
