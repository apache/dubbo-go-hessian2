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
	"reflect"
)

import (
	jerrors "github.com/juju/errors"
)

// hessian decode respone
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
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == byte(0x00) {
		return jerrors.Errorf("serialization ID:%v", header.SerialID)
	}

	flag := buf[2] & FLAG_EVENT
	if flag != byte(0x00) {
		header.Type |= Heartbeat
	}
	flag = buf[2] & FLAG_TWOWAY
	if flag != byte(0x00) {
		header.Type |= Response
	}
	flag = buf[2] & FLAG_REQUEST
	if flag != byte(0x00) {
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
	case RESPONSE_WITH_EXCEPTION:
		expt, err := decoder.Decode()
		if err != nil {
			return jerrors.Trace(err)
		}
		return jerrors.Errorf("got exception: %+v", expt)

	case RESPONSE_VALUE:
		rsp, err := decoder.Decode()
		if err != nil {
			return jerrors.Trace(err)
		}
		return jerrors.Trace(ReflectResponse(rsp, rspObj))

	case RESPONSE_NULL_VALUE:
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
