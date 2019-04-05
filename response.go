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

package hessian

import (
	"encoding/binary"
	"math"
	"reflect"
	"strconv"
	"strings"
)

import (
	jerrors "github.com/juju/errors"
)

// dubbo-remoting/dubbo-remoting-api/src/main/java/com/alibaba/dubbo/remoting/exchange/codec/ExchangeCodec.java
// v2.7.1 line 256 encodeResponse
// hessian encode response
func PackResponse(header DubboHeader, attachments map[string]string, ret interface{}) ([]byte, error) {
	var (
		byteArray []byte
	)

	hb := header.Type == Heartbeat

	// magic
	if hb {
		byteArray = append(byteArray, DubboResponseHeartbeatHeader[:]...)
	} else {
		byteArray = append(byteArray, DubboResponseHeaderBytes[:]...)
	}
	// set serialID, identify serialization types, eg: fastjson->6, hessian2->2
	byteArray[2] |= byte(header.SerialID & SERIAL_MASK)
	// response status
	if header.ResponseStatus != 0 {
		byteArray[3] = header.ResponseStatus
	}

	// request id
	binary.BigEndian.PutUint64(byteArray[4:], uint64(header.ID))

	// body
	encoder := NewEncoder()
	encoder.Append(byteArray[:HEADER_LENGTH])

	if hb {
		encoder.Encode(nil)
	} else {
		// com.alibaba.dubbo.rpc.protocol.dubbo.DubboCodec.DubboCodec.java
		// v2.7.1 line191 encodeRequestData

		atta := isSupportResponseAttachment(attachments[DUBBO_VERSION_KEY])

		var resWithException, resValue, resNullValue int32
		if atta {
			resWithException = RESPONSE_WITH_EXCEPTION
			resValue = RESPONSE_VALUE
			resNullValue = RESPONSE_NULL_VALUE
		} else {
			resWithException = RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS
			resValue = RESPONSE_VALUE_WITH_ATTACHMENTS
			resNullValue = RESPONSE_NULL_VALUE_WITH_ATTACHMENTS
		}

		if e, ok := ret.(error); ok { // throw error
			encoder.Encode(resWithException)
			encoder.Encode(e.Error())
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

	byteArray = encoder.Buffer()
	byteArray = encNull(byteArray) // if not, "java client" will throw exception  "unexpected end of file"
	pkgLen := len(byteArray)
	if pkgLen > int(DEFAULT_LEN) { // 8M
		return nil, jerrors.Errorf("Data length %d too large, max payload %d", pkgLen, DEFAULT_LEN)
	}
	// byteArray{body length}
	binary.BigEndian.PutUint32(byteArray[12:], uint32(pkgLen-HEADER_LENGTH))
	return byteArray, nil

}

// hessian decode response
func UnpackResponseHeaer(buf []byte, header *DubboHeader) error {
	// length := len(buf)
	// // hessianCodec.ReadHeader has check the header length
	// if length < HEADER_LENGTH {
	// 	return ErrHeaderNotEnough
	// }

	if buf[0] != byte(MAGIC_HIGH) && buf[1] != byte(MAGIC_LOW) {
		return ErrIllegalPackage
	}

	// Header{serialization id(5 bit), event, two way, req/response}
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == Zero {
		return jerrors.Errorf("serialization ID:%v", header.SerialID)
	}

	flag := buf[2] & FLAG_EVENT
	if flag != Zero {
		header.Type |= Heartbeat
	}
	flag = buf[3]
	if flag != Zero {
		header.Type |= Response
		header.ResponseStatus = flag
	}
	flag = buf[2] & FLAG_REQUEST
	if flag != Zero {
		return jerrors.Errorf("response flag:%v", flag)
	}

	// Header{status}
	var err error
	if buf[3] != Response_OK {
		err = ErrJavaException
		// return jerrors.Errorf("Response not OK, java exception:%s", string(buf[18:length-1]))
	}

	// Header{req id}
	header.ID = int64(binary.BigEndian.Uint64(buf[4:]))

	// Header{body len}
	header.BodyLen = int(binary.BigEndian.Uint32(buf[12:]))
	if header.BodyLen < 0 {
		return ErrIllegalPackage
	}

	return err
}

// hessian decode response body
func UnpackResponseBody(buf []byte, rspObj interface{}) error {
	// body
	decoder := NewDecoder(buf[:])
	rspType, err := decoder.Decode()
	if err != nil {
		return jerrors.Trace(err)
	}

	switch rspType {
	case RESPONSE_WITH_EXCEPTION, RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS:
		expt, err := decoder.Decode()
		if err != nil {
			return jerrors.Trace(err)
		}
		return jerrors.Errorf("got exception: %+v", expt)

	case RESPONSE_VALUE, RESPONSE_VALUE_WITH_ATTACHMENTS:
		rsp, err := decoder.Decode()
		if err != nil {
			return jerrors.Trace(err)
		}
		return jerrors.Trace(ReflectResponse(rsp, rspObj))

	case RESPONSE_NULL_VALUE, RESPONSE_NULL_VALUE_WITH_ATTACHMENTS:
		return jerrors.New("Received null")
	}

	return nil
}

func CopySlice(inSlice, outSlice reflect.Value) error {
	if inSlice.IsNil() {
		return jerrors.New("@in is nil")
	}
	if inSlice.Kind() != reflect.Slice {
		return jerrors.Errorf("@in is not slice, but %v", inSlice.Kind())
	}

	for outSlice.Kind() == reflect.Ptr {
		outSlice = outSlice.Elem()
	}

	size := inSlice.Len()
	outSlice.Set(reflect.MakeSlice(outSlice.Type(), size, size))

	for i := 0; i < size; i++ {
		inSliceValue := inSlice.Index(i)
		if !inSliceValue.Type().AssignableTo(outSlice.Index(i).Type()) {
			return jerrors.Errorf("in element type [%s] can not assign to out element type [%s]",
				inSliceValue.Type().String(), outSlice.Type().String())
		}
		outSlice.Index(i).Set(inSliceValue)
	}

	return nil
}

func CopyMap(inMapValue, outMapValue reflect.Value) error {
	if inMapValue.IsNil() {
		return jerrors.New("@in is nil")
	}
	if !inMapValue.CanInterface() {
		return jerrors.New("@in's Interface can not be used.")
	}
	if inMapValue.Kind() != reflect.Map {
		return jerrors.Errorf("@in is not map, but %v", inMapValue.Kind())
	}

	outMapType := UnpackPtrType(outMapValue.Type())
	SetValue(outMapValue, reflect.MakeMap(outMapType))

	outKeyType := outMapType.Key()

	outMapValue = UnpackPtrValue(outMapValue)
	outValueType := outMapValue.Type().Elem()

	for _, inKey := range inMapValue.MapKeys() {
		inValue := inMapValue.MapIndex(inKey)

		if !inKey.Type().AssignableTo(outKeyType) {
			return jerrors.Errorf("in Key:{type:%s, value:%#v} can not assign to out Key:{type:%s} ",
				inKey.Type().String(), inKey, outKeyType.String())
		}
		if !inValue.Type().AssignableTo(outValueType) {
			return jerrors.Errorf("in Value:{type:%s, value:%#v} can not assign to out value:{type:%s}",
				inValue.Type().String(), inValue, outValueType.String())
		}
		outMapValue.SetMapIndex(inKey, inValue)
	}

	return nil
}

// reflect return value
func ReflectResponse(in interface{}, out interface{}) error {
	if in == nil {
		return jerrors.Errorf("@in is nil")
	}

	if out == nil {
		return jerrors.Errorf("@out is nil")
	}
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return jerrors.Errorf("@out should be a pointer")
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

	if v >= 2001000 && v <= 2060200 {
		return false
	}
	return v >= LOWEST_VERSION_FOR_RESPONSE_ATTACHMENT
}

func version2Int(version string) int {
	var v = 0
	varr := strings.Split(version, ".")
	len := len(varr)
	for key, value := range varr {
		v0, err := strconv.Atoi(value)
		if err != nil {
			return -1
		}
		v += v0 * int(math.Pow10((len-key-1)*2))
	}
	return v
}
