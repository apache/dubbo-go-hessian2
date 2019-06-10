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
	"encoding/binary"
	"reflect"
	"strconv"
	"strings"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

/////////////////////////////////////////
// dubbo
/////////////////////////////////////////

// com.alibaba.dubbo.common.utils.ReflectUtils.ReflectUtils.java line245 getDesc
func getArgType(v interface{}) string {
	if v == nil {
		return "V"
	}

	switch v.(type) {
	// Serialized tags for base types
	case nil:
		return "V"
	case bool:
		return "Z"
	case byte:
		return "B"
	case int8:
		return "B"
	case int16:
		return "S"
	case uint16: // Equivalent to Char of Java
		return "C"
	// case rune:
	//	return "C"
	case int:
		return "I"
	case int32:
		return "I"
	case int64:
		return "J"
	case time.Time:
		return "java.util.Date"
	case float32:
		return "F"
	case float64:
		return "D"
	case string:
		return "java.lang.String"
	case []byte:
		return "[B"
	case map[interface{}]interface{}:
		// return  "java.util.HashMap"
		return "java.util.Map"

	//  Serialized tags for complex types
	default:
		t := reflect.TypeOf(v)
		if reflect.Ptr == t.Kind() {
			t = reflect.TypeOf(reflect.ValueOf(v).Elem())
		}
		switch t.Kind() {
		case reflect.Struct:
			return "java.lang.Object"
		case reflect.Slice, reflect.Array:
			// return "java.util.ArrayList"
			return "java.util.List"
		case reflect.Map: // Enter here, map may be map[string]int
			return "java.util.Map"
		default:
			return ""
		}
	}

	// unreachable
	// return "java.lang.RuntimeException"
}

func getArgsTypeList(args []interface{}) (string, error) {
	var (
		typ   string
		types string
	)

	for i := range args {
		typ = getArgType(args[i])
		if typ == "" {
			return types, perrors.Errorf("cat not get arg %#v type", args[i])
		}
		if !strings.Contains(typ, ".") {
			types += typ
		} else {
			// java.util.List -> Ljava/util/List;
			types += "L" + strings.Replace(typ, ".", "/", -1) + ";"
		}
	}

	return types, nil
}

// dubbo-remoting/dubbo-remoting-api/src/main/java/com/alibaba/dubbo/remoting/exchange/codec/ExchangeCodec.java
// v2.5.4 line 204 encodeRequest
// todo: attachments
func packRequest(service Service, header DubboHeader, params interface{}) ([]byte, error) {
	var (
		err           error
		types         string
		byteArray     []byte
		version       string
		pkgLen        int
		serviceParams map[string]string
	)

	args, ok := params.([]interface{})
	if !ok {
		return nil, perrors.Errorf("@params is not of type: []interface{}")
	}

	hb := header.Type == PackageHeartbeat

	//////////////////////////////////////////
	// byteArray
	//////////////////////////////////////////
	// magic
	switch header.Type {
	case PackageHeartbeat:
		byteArray = append(byteArray, DubboRequestHeartbeatHeader[:]...)
	case PackageRequest_TwoWay:
		byteArray = append(byteArray, DubboRequestHeaderBytesTwoWay[:]...)
	default:
		byteArray = append(byteArray, DubboRequestHeaderBytes[:]...)
	}

	// serialization id, two way flag, event, request/response flag
	// SerialID is id of serialization approach in java dubbo
	byteArray[2] |= header.GetSerialID()
	// request id
	binary.BigEndian.PutUint64(byteArray[4:], header.ID)

	encoder := NewEncoder()
	encoder.Append(byteArray[:HEADER_LENGTH])

	// com.alibaba.dubbo.rpc.protocol.dubbo.DubboCodec.DubboCodec.java line144 encodeRequestData
	//////////////////////////////////////////
	// body
	//////////////////////////////////////////
	if hb {
		encoder.Encode(nil)
		goto END
	}

	// dubbo version + path + version + method
	encoder.Encode(DUBBO_VERSION)
	encoder.Encode(service.Target)
	encoder.Encode(service.Version)
	encoder.Encode(service.Method)

	// args = args type list + args value list
	if types, err = getArgsTypeList(args); err != nil {
		return nil, perrors.Wrapf(err, " PackRequest(args:%+v)", args)
	}
	encoder.Encode(types)
	for _, v := range args {
		encoder.Encode(v)
	}

	serviceParams = make(map[string]string)
	serviceParams[PATH_KEY] = service.Path
	serviceParams[INTERFACE_KEY] = service.Interface
	if len(version) != 0 {
		serviceParams[VERSION_KEY] = version
	}
	if service.Timeout != 0 {
		serviceParams[TIMEOUT_KEY] = strconv.Itoa(int(service.Timeout / time.Millisecond))
	}

	encoder.Encode(serviceParams)

END:
	byteArray = encoder.Buffer()
	pkgLen = len(byteArray)
	if pkgLen > int(DEFAULT_LEN) { // 8M
		return nil, perrors.Errorf("Data length %d too large, max payload %d", pkgLen, DEFAULT_LEN)
	}
	// byteArray{body length}
	binary.BigEndian.PutUint32(byteArray[12:], uint32(pkgLen-HEADER_LENGTH))
	return byteArray, nil
}

// hessian decode request body
func unpackRequestBody(buf []byte, reqObj interface{}) error {

	req, ok := reqObj.([]interface{})
	if !ok {
		return perrors.Errorf("@reqObj is not of type: []interface{}")
	}
	if len(req) < 7 {
		return perrors.New("length of @reqObj should  be 7")
	}

	var (
		err                                                     error
		dubboVersion, target, serviceVersion, method, argsTypes interface{}
		args                                                    []interface{}
	)
	decoder := NewDecoder(buf[:])

	dubboVersion, err = decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[0] = dubboVersion

	target, err = decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[1] = target

	serviceVersion, err = decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[2] = serviceVersion

	method, err = decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[3] = method

	argsTypes, err = decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[4] = argsTypes

	ats := DescRegex.FindAllString(argsTypes.(string), -1)
	var arg interface{}
	for i := 0; i < len(ats); i++ {
		arg, err = decoder.Decode()
		if err != nil {
			return perrors.WithStack(err)
		}
		args = append(args, arg)
	}
	req[5] = args

	attachments, err := decoder.Decode()
	if err != nil {
		return perrors.WithStack(err)
	}
	req[6] = attachments

	return nil
}
