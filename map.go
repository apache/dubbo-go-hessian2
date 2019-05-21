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
	"io"
	"reflect"
)

import (
	jerrors "github.com/juju/errors"
)

/////////////////////////////////////////
// map/object
/////////////////////////////////////////

// ::= 'M' type (value value)* 'Z'  # key, value map pairs
// ::= 'H' (value value)* 'Z'       # untyped key, value
func (e *Encoder) encUntypedMap(m map[interface{}]interface{}) error {
	if len(m) == 0 {
		return nil
	}

	// check ref
	if n, ok := e.checkRefMap(reflect.ValueOf(m)); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	var err error
	e.buffer = encByte(e.buffer, BC_MAP_UNTYPED)
	for k, v := range m {
		if err = e.Encode(k); err != nil {
			return err
		}
		if err = e.Encode(v); err != nil {
			return err
		}
	}

	e.buffer = encByte(e.buffer, BC_END) // 'Z'

	return nil
}

func getMapKey(key reflect.Value, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Bool:
		return key.Bool(), nil

	case reflect.Int8:
		return int8(key.Int()), nil
	case reflect.Int16:
		return int16(key.Int()), nil
	case reflect.Int32:
		return int32(key.Int()), nil
	case reflect.Int:
		return int(key.Int()), nil
	case reflect.Int64:
		return key.Int(), nil

	case reflect.Uint8:
		return byte(key.Uint()), nil
	case reflect.Uint16:
		return uint16(key.Uint()), nil
	case reflect.Uint32:
		return uint32(key.Uint()), nil
	case reflect.Uint:
		return uint(key.Uint()), nil
	case reflect.Uint64:
		return key.Uint(), nil

	case reflect.Float32:
		return float32(key.Float()), nil
	case reflect.Float64:
		return key.Float(), nil

	case reflect.Uintptr:
		return key.UnsafeAddr(), nil

	case reflect.String:
		return key.String(), nil
	}

	return nil, jerrors.Errorf("unsupported map key kind %s", t.Kind().String())
}

func (e *Encoder) encMap(m interface{}) error {
	var (
		err   error
		k     interface{}
		typ   reflect.Type
		value reflect.Value
		keys  []reflect.Value
	)

	value = reflect.ValueOf(m)

	// check ref
	if n, ok := e.checkRefMap(value); ok {
		e.buffer = encRef(e.buffer, n)
		return nil
	}

	value = UnpackPtrValue(value)
	// check nil map
	if value.Kind() == reflect.Ptr && !value.Elem().IsValid() {
		e.buffer = encNull(e.buffer)
		return nil
	}

	keys = value.MapKeys()
	if len(keys) == 0 {
		// fix: set nil for empty map
		e.buffer = encNull(e.buffer)
		return nil
	}

	typ = value.Type().Key()
	e.buffer = encByte(e.buffer, BC_MAP_UNTYPED)
	for i := 0; i < len(keys); i++ {
		k, err = getMapKey(keys[i], typ)
		if err != nil {
			return jerrors.Annotatef(err, "getMapKey(idx:%d, key:%+v)", i, keys[i])
		}
		if err = e.Encode(k); err != nil {
			return jerrors.Annotatef(err, "failed to encode map key(idx:%d, key:%+v)", i, keys[i])
		}
		entryValueObj := value.MapIndex(keys[i]).Interface()
		if err = e.Encode(entryValueObj); err != nil {
			return jerrors.Annotatef(err, "failed to encode map value(idx:%d, key:%+v, value:%+v)", i, k, entryValueObj)
		}
	}
	e.buffer = encByte(e.buffer, BC_END)

	return nil
}

/////////////////////////////////////////
// Map
/////////////////////////////////////////

// ::= 'M' type (value value)* 'Z'  # key, value map pairs
// ::= 'H' (value value)* 'Z'       # untyped key, value
func (d *Decoder) decMapByValue(value reflect.Value) error {
	var (
		tag        byte
		err        error
		entryKey   interface{}
		entryValue interface{}
	)

	//tag, _ = d.readBufByte()
	tag, err = d.readByte()
	// check error
	if err != nil {
		return jerrors.Trace(err)
	}

	switch tag {
	case BC_NULL:
		// null map tag check
		return nil
	case BC_REF:
		refObj, err := d.decRef(int32(tag))
		if err != nil {
			return jerrors.Trace(err)
		}
		SetValue(value, EnsurePackValue(refObj))
		return nil
	case BC_MAP:
		d.decString(TAG_READ) // read map type , ignored
	case BC_MAP_UNTYPED:
		//do nothing
	default:
		return jerrors.Errorf("expect map header, but get %x", tag)
	}

	m := reflect.MakeMap(UnpackPtrType(value.Type()))
	// pack with pointer, so that to ref the same map
	m = PackPtr(m)
	d.appendRefs(m)

	//read key and value
	for {
		entryKey, err = d.Decode()
		if err != nil {
			// EOF means the end flag 'Z' of map is already read
			if err == io.EOF {
				break
			} else {
				return jerrors.Trace(err)
			}
		}
		if entryKey == nil {
			break
		}
		entryValue, err = d.Decode()
		// fix: check error
		if err != nil {
			return jerrors.Trace(err)
		}
		m.Elem().SetMapIndex(EnsurePackValue(entryKey), EnsurePackValue(entryValue))
	}

	SetValue(value, m)

	return nil
}

func (d *Decoder) decMap(flag int32) (interface{}, error) {
	var (
		err        error
		tag        byte
		ok         bool
		k          interface{}
		v          interface{}
		t          string
		keyName    string
		methodName string
		key        interface{}
		value      interface{}
		inst       interface{}
		m          map[interface{}]interface{}
		fieldValue reflect.Value
		args       []reflect.Value
	)

	if flag != TAG_READ {
		tag = byte(flag)
	} else {
		tag, _ = d.readByte()
	}

	switch {
	case tag == BC_NULL:
		return nil, nil
	case tag == BC_REF:
		return d.decRef(int32(tag))
	case tag == BC_MAP:
		if t, err = d.decType(); err != nil {
			return nil, err
		}

		if _, ok = checkPOJORegistry(t); ok {
			m = make(map[interface{}]interface{}) // todo: This assumes the definition of map, which is incorrect.
			d.appendRefs(m)

			// d.decType() // ignore
			for d.peekByte() != byte('z') {
				k, err = d.Decode()
				if err != nil {
					if err == io.EOF {
						break
					}

					return nil, err
				}
				v, err = d.Decode()
				if err != nil {
					return nil, err
				}
				m[k] = v
			}
			_, err = d.readByte()
			// check error
			if err != nil {
				return nil, jerrors.Trace(err)
			}

			return m, nil
		}

		// check failed
		inst = createInstance(t)
		d.appendRefs(inst)

		for d.peekByte() != 'z' {
			if key, err = d.Decode(); err != nil {
				return nil, err
			}
			if value, err = d.Decode(); err != nil {
				return nil, err
			}
			//set value of the struct to Zero
			if fieldValue = reflect.ValueOf(value); fieldValue.IsValid() {
				keyName = key.(string)
				if keyName[0] >= 'a' { //convert to Upper
					methodName = "Set" + string(keyName[0]-32) + keyName[1:]
				} else {
					methodName = "Set" + keyName
				}

				args = args[:0]
				args = append(args, fieldValue)
				reflect.ValueOf(inst).MethodByName(methodName).Call(args)
			}
		}

		return inst, nil

	case tag == BC_MAP_UNTYPED:
		m = make(map[interface{}]interface{})
		d.appendRefs(m)
		for d.peekByte() != BC_END {
			k, err = d.Decode()
			if err != nil {
				if err == io.EOF {
					break
				}

				return nil, err
			}
			v, err = d.Decode()
			if err != nil {
				return nil, err
			}
			m[k] = v
		}
		_, err = d.readByte()
		// check error
		if err != nil {
			return nil, jerrors.Trace(err)
		}
		return m, nil

	default:
		return nil, jerrors.Errorf("illegal map type tag:%+v", tag)
	}
}
