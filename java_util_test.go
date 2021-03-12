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
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_util"
)

func TestJavaUtil(t *testing.T) {
	testDecodeFramework(t, "javautilUUID", &java_util.UUID{LeastSigBits: int64(-7160773830801198154), MostSigBits: int64(459021424248441700)})
}

// TestJavaRandomUUID is test java UUID.toString() equals go UUID.String()
// java test result include uuid encode and uuid.toString
// use '@@@' split that
func TestJavaRandomUUID(t *testing.T) {
	method := "javautilRandomUUID"
	b := getJavaReply(method, "")
	split := strings.Split(string(b), "@@@")
	d := NewDecoder([]byte(split[0]))
	r, e := d.Decode()
	if e != nil {
		t.Errorf(e.Error())
	}
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}
	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	uuid, ok := r.(*java_util.UUID)
	if ok {
		ok := reflect.DeepEqual(split[1], uuid.String())
		if !ok {
			t.Error("go String() no equal java toString()")
		}
	} else {
		t.Error("decode UUID struct false")
	}
}
