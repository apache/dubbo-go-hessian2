// Copyright 2016-2019 Alex Stocks, Yincheng Fang
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
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
	"strconv"
	"strings"
)

import (
	perrors "github.com/pkg/errors"
)

type Response struct {
	RspObj    interface{}
	Exception error
	//Attachments map[string]string
}

// dubbo-remoting/dubbo-remoting-api/src/main/java/com/alibaba/dubbo/remoting/exchange/codec/ExchangeCodec.java
// v2.7.1 line 256 encodeResponse
// hessian encode response
func packResponse(packageType PackageType, header DubboHeader, attachments map[string]string, ret interface{}) ([]byte, error) {
	var (
		byteArray []byte
	)

	header.MagicNumber = MAGIC
	header.SetPackageType(false, packageType)

	buf := bytes.NewBuffer(make([]byte, 0, HEADER_LENGTH))
	binary.Write(buf, binary.BigEndian, header)
	byteArray = buf.Bytes()

	// body
	encoder := NewEncoder()
	encoder.Append(byteArray)

	if header.ResponseStatus == Response_OK {
		hb := packageType == PackageHeartbeat
		if hb {
			encoder.Encode(nil)
		} else {
			// com.alibaba.dubbo.rpc.protocol.dubbo.DubboCodec.DubboCodec.java
			// v2.7.1 line191 encodeResponseData

			atta := isSupportResponseAttachment(attachments[DUBBO_VERSION_KEY])

			var resWithException, resValue, resNullValue int32
			if atta {
				resWithException = RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS
				resValue = RESPONSE_VALUE_WITH_ATTACHMENTS
				resNullValue = RESPONSE_NULL_VALUE_WITH_ATTACHMENTS
			} else {
				resWithException = RESPONSE_WITH_EXCEPTION
				resValue = RESPONSE_VALUE
				resNullValue = RESPONSE_NULL_VALUE
			}

			if e, ok := ret.(error); ok { // throw error
				encoder.Encode(resWithException)
				if t, ok := e.(Throwabler); ok {
					encoder.Encode(t)
				} else {
					encoder.Encode(e.Error())
				}
			} else {
				if ret == nil {
					encoder.Encode(resNullValue)
				} else {
					encoder.Encode(resValue)
					encoder.Encode(ret) // result
				}
			}

			if atta {
				encoder.Encode(attachments) // attachments
			}
		}
	} else {
		// com.alibaba.dubbo.remoting.exchange.codec.ExchangeCodec
		// v2.6.5 line280 encodeResponse
		if e, ok := ret.(error); ok { // throw error
			encoder.Encode(e.Error())
		} else if e, ok := ret.(string); ok {
			encoder.Encode(e)
		} else {
			return nil, perrors.New("Ret must be error or string!")
		}
	}

	byteArray = encoder.Buffer()
	byteArray = encNull(byteArray) // if not, "java client" will throw exception  "unexpected end of file"
	pkgLen := len(byteArray)
	if pkgLen > int(DEFAULT_LEN) { // 8M
		return nil, perrors.Errorf("Data length %d too large, max payload %d", pkgLen, DEFAULT_LEN)
	}
	// byteArray{body length}
	binary.BigEndian.PutUint32(byteArray[12:], uint32(pkgLen-HEADER_LENGTH))
	return byteArray, nil

}

// hessian decode response body
// todo: need to read attachments
func unpackResponseBody(buf []byte, response *Response) error {
	// body
	decoder := NewDecoder(buf[:])
	rspType, err := decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}

	switch rspType {
	case RESPONSE_WITH_EXCEPTION, RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS:
		expt, err := decoder.Decode()
		if err != nil {
			return perrors.WithStack(err)
		}
		if e, ok := expt.(error); ok {
			response.Exception = e
			return nil
		}
		response.Exception = perrors.Errorf("got exception: %+v", expt)
		return nil

	case RESPONSE_VALUE, RESPONSE_VALUE_WITH_ATTACHMENTS:
		rsp, err := decoder.Decode()
		if err != nil {
			return perrors.WithStack(err)
		}
		return perrors.WithStack(ReflectResponse(rsp, response.RspObj))

	case RESPONSE_NULL_VALUE, RESPONSE_NULL_VALUE_WITH_ATTACHMENTS:
		return nil
	}

	return nil
}

// CopySlice copy from inSlice to outSlice
func CopySlice(inSlice, outSlice reflect.Value) error {
	if inSlice.IsNil() {
		return perrors.New("@in is nil")
	}
	if inSlice.Kind() != reflect.Slice {
		return perrors.Errorf("@in is not slice, but %v", inSlice.Kind())
	}

	for outSlice.Kind() == reflect.Ptr {
		outSlice = outSlice.Elem()
	}

	size := inSlice.Len()
	outSlice.Set(reflect.MakeSlice(outSlice.Type(), size, size))

	for i := 0; i < size; i++ {
		inSliceValue := inSlice.Index(i)
		if !inSliceValue.Type().AssignableTo(outSlice.Index(i).Type()) {
			return perrors.Errorf("in element type [%s] can not assign to out element type [%s]",
				inSliceValue.Type().String(), outSlice.Type().String())
		}
		outSlice.Index(i).Set(inSliceValue)
	}

	return nil
}

// CopyMap copy from in map to out map
func CopyMap(inMapValue, outMapValue reflect.Value) error {
	if inMapValue.IsNil() {
		return perrors.New("@in is nil")
	}
	if !inMapValue.CanInterface() {
		return perrors.New("@in's Interface can not be used.")
	}
	if inMapValue.Kind() != reflect.Map {
		return perrors.Errorf("@in is not map, but %v", inMapValue.Kind())
	}

	outMapType := UnpackPtrType(outMapValue.Type())
	SetValue(outMapValue, reflect.MakeMap(outMapType))

	outKeyType := outMapType.Key()

	outMapValue = UnpackPtrValue(outMapValue)
	outValueType := outMapValue.Type().Elem()

	for _, inKey := range inMapValue.MapKeys() {
		inValue := inMapValue.MapIndex(inKey)

		if !inKey.Type().AssignableTo(outKeyType) {
			return perrors.Errorf("in Key:{type:%s, value:%#v} can not assign to out Key:{type:%s} ",
				inKey.Type().String(), inKey, outKeyType.String())
		}
		if !inValue.Type().AssignableTo(outValueType) {
			return perrors.Errorf("in Value:{type:%s, value:%#v} can not assign to out value:{type:%s}",
				inValue.Type().String(), inValue, outValueType.String())
		}
		outMapValue.SetMapIndex(inKey, inValue)
	}

	return nil
}

// ReflectResponse reflect return value
// TODO response object should not be copied again to another object, it should be the exact type of the object
func ReflectResponse(in interface{}, out interface{}) error {
	if in == nil {
		return perrors.Errorf("@in is nil")
	}

	if out == nil {
		return perrors.Errorf("@out is nil")
	}
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return perrors.Errorf("@out should be a pointer")
	}

	inValue := EnsurePackValue(in)
	outValue := EnsurePackValue(out)

	switch inValue.Type().Kind() {
	case reflect.Slice, reflect.Array:
		return CopySlice(inValue, outValue)
	case reflect.Map:
		return CopyMap(inValue, outValue)
	default:
		SetValue(outValue, inValue)
	}

	return nil
}

var versionInt = make(map[string]int)

// isSupportResponseAttachment is for compatibility among some dubbo version
// but we haven't used it yet.
// dubbo-common/src/main/java/org/apache/dubbo/common/Version.java
// v2.7.1 line 96
func isSupportResponseAttachment(version string) bool {
	if version == "" {
		return false
	}

	v, ok := versionInt[version]
	if !ok {
		v = version2Int(version)
		if v == -1 {
			return false
		}
	}

	if v >= 2001000 && v <= 2060200 { // 2.0.10 ~ 2.6.2
		return false
	}
	return v >= LOWEST_VERSION_FOR_RESPONSE_ATTACHMENT
}

func version2Int(version string) int {
	var v = 0
	varr := strings.Split(version, ".")
	length := len(varr)
	for key, value := range varr {
		v0, err := strconv.Atoi(value)
		if err != nil {
			return -1
		}
		v += v0 * int(math.Pow10((length-key-1)*2))
	}
	if length == 3 {
		return v * 100
	}
	return v
}
