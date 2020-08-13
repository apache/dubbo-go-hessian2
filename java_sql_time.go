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
	"io"
	"reflect"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

func init() {
	RegisterPOJO(&Date{})
	RegisterPOJO(&Time{})
}

type JavaSqlTime interface {
	SetTime(time time.Time)
	ValueOf(time string) (JavaSqlTime, error)
	JavaClassName() string
	UnixNano() int64
	time() time.Time
}
type Date struct {
	time.Time
}

func (d Date) time() time.Time {
	return d.Time
}

func (Date) JavaClassName() string {
	return "java.sql.Date"
}

func (d *Date) SetTime(time time.Time) {
	d.Time = time
}
func (Date) ValueOf(time string) (JavaSqlTime, error) {
	panic("")
}

type Time struct {
	time.Time
}

func (Time) JavaClassName() string {
	return "java.sql.Time"
}

func (t Time) time() time.Time {
	return t.Time
}

func (t *Time) SetTime(time time.Time) {
	t.Time = time
}
func (Time) ValueOf(time string) (JavaSqlTime, error) {
	panic("")
}

var javaSqlTimeTypeMap = make(map[string]reflect.Type, 16)

func SetJavaSqlTimeSerialize(time JavaSqlTime) {
	name := time.JavaClassName()
	v := reflect.ValueOf(time)
	var typ reflect.Type
	switch v.Kind() {
	case reflect.Struct:
		typ = v.Type()
	case reflect.Ptr:
		typ = v.Elem().Type()
	default:
		typ = reflect.TypeOf(time)
	}
	SetSerializer(name, JavaSqlTimeSerializer{})
	//RegisterPOJO(time)
	javaSqlTimeTypeMap[name] = typ
}

func getJavaSqlTimeSerialize(name string) reflect.Type {
	return javaSqlTimeTypeMap[name]
}

type JavaSqlTimeSerializer struct {
}

func (JavaSqlTimeSerializer) EncObject(e *Encoder, vv POJO) error {

	var (
		idx    int
		idx1   int
		i      int
		err    error
		clsDef classInfo
	)
	v, ok := vv.(JavaSqlTime)
	if !ok {
		return perrors.New("can not be converted into java sql time object")
	}
	className := v.JavaClassName()
	if className == "" {
		return perrors.New("class name empty")
	}

	tValue := reflect.ValueOf(vv)
	// check ref
	if n, ok := e.checkRefMap(tValue); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	// write object definition
	idx = -1
	for i = range e.classInfoList {
		if v.JavaClassName() == e.classInfoList[i].javaName {
			idx = i
			break
		}
	}

	if idx == -1 {
		idx1, ok = checkPOJORegistry(typeof(v))
		if !ok {
			if reflect.TypeOf(v).Implements(javaEnumType) {
				idx1 = RegisterJavaEnum(v.(POJOEnum))
			} else {
				idx1 = RegisterPOJO(v)
			}
		}
		_, clsDef, err = getStructDefByIndex(idx1)
		if err != nil {
			return perrors.WithStack(err)
		}

		i = len(e.classInfoList)
		e.classInfoList = append(e.classInfoList, clsDef)
		e.buffer = append(e.buffer, clsDef.buffer...)
		e.buffer = e.buffer[0 : len(e.buffer)-1]
	}

	//if idx < -1 {
	//	e.buffer = encString(e.buffer, "value")
	//	var bytes []byte
	//	bytes = append(bytes, PackInt64(v.UnixNano()/1e6)...)
	//	e.buffer = encDateInMs(bytes, v)
	//	e.buffer = append(e.buffer, BC_END)
	//} else
	//{
	if idx == -1 {
		e.buffer = encInt32(e.buffer, 1)
		e.buffer = encString(e.buffer, "value")

		// write object instance
		if byte(i) <= OBJECT_DIRECT_MAX {
			e.buffer = encByte(e.buffer, byte(i)+BC_OBJECT_DIRECT)
		} else {
			e.buffer = encByte(e.buffer, BC_OBJECT)
			e.buffer = encInt32(e.buffer, int32(idx1))
		}
		//if (ref == -1) {
		//	out.writeInt(1);
		//	out.writeString("value");
		//	out.writeObjectBegin(cl.getName());
		//}
		//}
		//var bytes []byte
		//bytes = append(bytes, PackInt64(v.UnixNano()/1e6)...)
		e.buffer = encDateInMs(e.buffer, v.time())
	}

	return nil
}

func (JavaSqlTimeSerializer) DecObject(d *Decoder, typ reflect.Type, cls classInfo) (interface{}, error) {

	if typ.Kind() != reflect.Struct {
		return nil, perrors.Errorf("wrong type expect Struct but get:%s", typ.String())
	}

	vRef := reflect.New(typ)
	// add pointer ref so that ref the same object
	d.appendRefs(vRef.Interface())

	tag, err := d.readByte()
	if err == io.EOF {
		return nil, err
	}
	date, err := d.decDate(int32(tag))
	if err != nil {
		date = date
	}
	sqlTime := vRef.Interface()

	result, ok := sqlTime.(JavaSqlTime)
	result.SetTime(date)
	if !ok {
		panic("result type is not sql time, please check the whether the conversion is ok")
	}
	return result, nil
}
