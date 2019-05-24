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
	"time"
	"unsafe"
)

import (
	perrors "github.com/pkg/errors"
)

// nil bool int8 int32 int64 float32 float64 time.Time
// string []byte []interface{} map[interface{}]interface{}
// array object struct

// Encoder struct
type Encoder struct {
	classInfoList []classInfo
	buffer        []byte
	refMap        map[unsafe.Pointer]_refElem
}

// NewEncoder generate an encoder instance
func NewEncoder() *Encoder {
	var buffer = make([]byte, 64)

	return &Encoder{
		buffer: buffer[:0],
		refMap: make(map[unsafe.Pointer]_refElem, 7),
	}
}

// Buffer returns byte buffer
func (e *Encoder) Buffer() []byte {
	return e.buffer[:]
}

// Append byte arr to encoder buffer
func (e *Encoder) Append(buf []byte) {
	e.buffer = append(e.buffer, buf[:]...)
}

// Encode If @v can not be encoded, the return value is nil. At present only struct may can not be encoded.
func (e *Encoder) Encode(v interface{}) error {
	if v == nil {
		e.buffer = encNull(e.buffer)
		return nil
	}

	switch val := v.(type) {
	case nil:
		e.buffer = encNull(e.buffer)
		return nil

	case bool:
		e.buffer = encBool(e.buffer, val)

	case int8, int16, int32:
		e.buffer = encInt32(e.buffer, v.(int32))

	case int:
		// if v.(int) >= -2147483648 && v.(int) <= 2147483647 {
		// 	b = encInt32(int32(v.(int)), b)
		// } else {
		// 	b = encInt64(int64(v.(int)), b)
		// }
		// use int64 type to handle int, to avoid  panic like :  reflect: Call using int32 as type int64 [recovered]
		// when decode
		e.buffer = encInt64(e.buffer, int64(val))

	case int64:
		e.buffer = encInt64(e.buffer, val)

	case time.Time:
		e.buffer = encDateInMs(e.buffer, val)
		// e.buffer = encDateInMimute(v.(time.Time), e.buffer)

	case float32:
		e.buffer = encFloat(e.buffer, float64(val))

	case float64:
		e.buffer = encFloat(e.buffer, val)

	case string:
		e.buffer = encString(e.buffer, val)

	case []byte:
		e.buffer = encBinary(e.buffer, val)

	case map[interface{}]interface{}:
		return e.encUntypedMap(val)

	default:
		t := UnpackPtrType(reflect.TypeOf(v))
		switch t.Kind() {
		case reflect.Struct:
			if p, ok := v.(POJO); ok {
				return e.encObject(p)
			}

			return perrors.Errorf("struct type not Support! %s[%v] is not a instance of POJO!", t.String(), v)
		case reflect.Slice, reflect.Array:
			return e.encUntypedList(v)
		case reflect.Map: // the type must be map[string]int
			return e.encMap(v)
		default:
			if p, ok := v.(POJOEnum); ok { // JavaEnum
				return e.encObject(p)
			}
			return perrors.Errorf("type not supported! %s", t.Kind().String())
		}
	}

	return nil
}
