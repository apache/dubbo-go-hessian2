// Copyright 2016-2019 Alex Stocks, Wongoo, Yincheng Fang
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
	"math"
	"reflect"
	"testing"
)

type Department struct {
	Name string
}

func (Department) JavaClassName() string {
	return "com.bdt.info.Department"
}

type WorkerInfo struct {
	unexportedFiled   string
	Name              string
	Addrress          string
	Age               int
	Salary            float32
	Payload           map[string]int32
	FamilyMembers     []string `hessian:"familyMembers1"`
	FamilyPhoneNumber string   // default convert to => familyPhoneNumber
	Dpt               Department
}

func (WorkerInfo) JavaClassName() string {
	return "com.bdt.info.WorkerInfo"
}

func TestEncEmptyStruct(t *testing.T) {
	var (
		w   WorkerInfo
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	w = WorkerInfo{
		Name:              "Trump",
		Addrress:          "W,D.C.",
		Age:               72,
		Salary:            21000.03,
		Payload:           map[string]int32{"Number": 2017061118},
		FamilyMembers:     []string{"m1", "m2", "m3"},
		FamilyPhoneNumber: "010-12345678",
		// Dpt: Department{
		// 	Name: "Adm",
		// },
	}
	e.Encode(w)

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", w, res, err)

	reflect.DeepEqual(w, res)
}

func TestEncStruct(t *testing.T) {
	var (
		w   WorkerInfo
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	w = WorkerInfo{
		Name:              "Trump",
		Addrress:          "W,D.C.",
		Age:               72,
		Salary:            21000.03,
		Payload:           map[string]int32{"Number": 2017061118},
		FamilyMembers:     []string{"m1", "m2", "m3"},
		FamilyPhoneNumber: "010-12345678",
		Dpt: Department{
			Name: "Adm",
		},
		unexportedFiled: "you cannot see me!",
	}
	wCopy := w
	wCopy.unexportedFiled = ""

	err = e.Encode(w)
	if err != nil {
		t.Errorf("Encode() = %+v", err)
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%+v) = %+v, %v\n", w, res, err)

	w2, ok := res.(*WorkerInfo)
	if !ok {
		t.Fatalf("res:%T is not of type WorkerInfo", w2)
	}

	if !reflect.DeepEqual(wCopy, *w2) {
		t.Fatalf("w != w2:\n%#v\n!=\n%#v", wCopy, w2)
	}
}

type UserName struct {
	FirstName string
	LastName  string
}

func (UserName) JavaClassName() string {
	return "com.bdt.info.UserName"
}

type Person struct {
	UserName
	Age int32
	Sex bool
}

func (Person) JavaClassName() string {
	return "com.bdt.info.Person"
}

type JOB struct {
	Title   string
	Company string
}

func (JOB) JavaClassName() string {
	return "com.bdt.info.JOB"
}

type Worker struct {
	Person
	CurJob JOB
	Jobs   []JOB
}

func (Worker) JavaClassName() string {
	return "com.bdt.info.Worker"
}

func TestIssue6(t *testing.T) {
	name := UserName{
		FirstName: "John",
		LastName:  "Doe",
	}
	person := Person{
		UserName: name,
		Age:      18,
		Sex:      true,
	}

	worker := &Worker{
		Person: person,
		CurJob: JOB{Title: "cto", Company: "facebook"},
		Jobs: []JOB{
			JOB{Title: "manager", Company: "google"},
			JOB{Title: "ceo", Company: "microsoft"},
		},
	}

	e := NewEncoder()
	err := e.Encode(worker)
	if err != nil {
		t.Fatalf("encode(worker:%#v) = error:%s", worker, err)
	}
	bytes := e.Buffer()

	d := NewDecoder(bytes)
	res, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("type of decode object:%v", reflect.TypeOf(res))

	worker2, ok := res.(*Worker)
	if !ok {
		t.Fatalf("res:%#v is not of type Worker", res)
	}

	if !reflect.DeepEqual(worker, worker2) {
		t.Fatalf("worker:%#v != worker2:%#v", worker, worker2)
	}
}

type A0 struct{}

func (*A0) JavaClassName() string {
	return "com.caucho.hessian.test.A0"
}

type A1 struct{}

func (*A1) JavaClassName() string {
	return "com.caucho.hessian.test.A1"
}

type A2 struct{}

func (*A2) JavaClassName() string {
	return "com.caucho.hessian.test.A2"
}

type A3 struct{}

func (*A3) JavaClassName() string {
	return "com.caucho.hessian.test.A3"
}

type A4 struct{}

func (*A4) JavaClassName() string {
	return "com.caucho.hessian.test.A4"
}

type A5 struct{}

func (*A5) JavaClassName() string {
	return "com.caucho.hessian.test.A5"
}

type A6 struct{}

func (*A6) JavaClassName() string {
	return "com.caucho.hessian.test.A6"
}

type A7 struct{}

func (*A7) JavaClassName() string {
	return "com.caucho.hessian.test.A7"
}

type A8 struct{}

func (*A8) JavaClassName() string {
	return "com.caucho.hessian.test.A8"
}

type A9 struct{}

func (*A9) JavaClassName() string {
	return "com.caucho.hessian.test.A9"
}

type A10 struct{}

func (*A10) JavaClassName() string {
	return "com.caucho.hessian.test.A10"
}

type A11 struct{}

func (*A11) JavaClassName() string {
	return "com.caucho.hessian.test.A11"
}

type A12 struct{}

func (*A12) JavaClassName() string {
	return "com.caucho.hessian.test.A12"
}

type A13 struct{}

func (*A13) JavaClassName() string {
	return "com.caucho.hessian.test.A13"
}

type A14 struct{}

func (*A14) JavaClassName() string {
	return "com.caucho.hessian.test.A14"
}

type A15 struct{}

func (*A15) JavaClassName() string {
	return "com.caucho.hessian.test.A15"
}

type A16 struct{}

func (*A16) JavaClassName() string {
	return "com.caucho.hessian.test.A16"
}

type TestObjectStruct struct {
	Value int32 `hessian:"_value"`
}

func (*TestObjectStruct) JavaClassName() string {
	return "com.caucho.hessian.test.TestObject"
}

type TestConsStruct struct {
	First string          `hessian:"_first"`
	Rest  *TestConsStruct `hessian:"_rest"`
}

func (*TestConsStruct) JavaClassName() string {
	return "com.caucho.hessian.test.TestCons"
}

func TestObject(t *testing.T) {
	RegisterPOJOs(
		&A0{},
		&A1{},
		&A2{},
		&A3{},
		&A4{},
		&A5{},
		&A6{},
		&A7{},
		&A8{},
		&A9{},
		&A10{},
		&A11{},
		&A12{},
		&A13{},
		&A14{},
		&A15{},
		&A16{},
		&TestObjectStruct{},
		&TestConsStruct{},
	)

	testDecodeFramework(t, "replyObject_0", &A0{})
	testDecodeFramework(t, "replyObject_1", &TestObjectStruct{Value: 0})
	testDecodeFramework(t, "replyObject_16", []interface{}{&A0{}, &A1{}, &A2{}, &A3{}, &A4{}, &A5{}, &A6{}, &A7{}, &A8{}, &A9{}, &A10{}, &A11{}, &A12{}, &A13{}, &A14{}, &A15{}, &A16{}})
	testDecodeFramework(t, "replyObject_2", []interface{}{&TestObjectStruct{Value: 0}, &TestObjectStruct{Value: 1}})
	testDecodeFramework(t, "replyObject_2b", []interface{}{&TestObjectStruct{Value: 0}, &TestObjectStruct{Value: 0}})

	object := TestObjectStruct{Value: 0}
	object2a := []interface{}{&object, &object}
	testDecodeFramework(t, "replyObject_2a", object2a)

	cons := TestConsStruct{}
	cons.First = "a"
	cons.Rest = &cons
	testDecodeFramework(t, "replyObject_3", &cons)
}

func TestObjectEncode(t *testing.T) {
	testJavaDecode(t, "argObject_0", &A0{})
	testJavaDecode(t, "argObject_1", &TestObjectStruct{Value: 0})
	testJavaDecode(t, "argObject_16", []interface{}{&A0{}, &A1{}, &A2{}, &A3{}, &A4{}, &A5{}, &A6{}, &A7{}, &A8{}, &A9{}, &A10{}, &A11{}, &A12{}, &A13{}, &A14{}, &A15{}, &A16{}})
	testJavaDecode(t, "argObject_2", []interface{}{&TestObjectStruct{Value: 0}, &TestObjectStruct{Value: 1}})
	testJavaDecode(t, "argObject_2b", []interface{}{&TestObjectStruct{Value: 0}, &TestObjectStruct{Value: 0}})

	object := TestObjectStruct{Value: 0}
	object2a := []interface{}{&object, &object}
	testJavaDecode(t, "argObject_2a", object2a)

	cons := TestConsStruct{}
	cons.First = "a"
	cons.Rest = &cons
	testJavaDecode(t, "argObject_3", &cons)
}

type Tuple struct {
	Byte    int8
	Short   int16
	Integer int32
	Long    int64
	Double  float32
	B       uint8
	S       uint16
	I       uint32
	L       uint64
	D       float64
}

func (t Tuple) JavaClassName() string {
	return "test.tuple.Tuple"
}

func TestDecodeJavaTupleObject(t *testing.T) {
	tuple := &Tuple{
		Byte:    1,
		Short:   1,
		Integer: 1,
		Long:    1,
		Double:  1.23,
		B:       0x01,
		S:       1,
		I:       1,
		L:       1,
		D:       1.23,
	}

	RegisterPOJO(tuple)

	testDecodeJavaData(t, "getTheTuple", "test.tuple.TupleProviderImpl", tuple)
}

func TestEncodeDecodeTuple(t *testing.T) {
	doTestEncodeDecodeTuple(t, &Tuple{
		Byte:    1,
		Short:   1,
		Integer: 1,
		Long:    1,
		Double:  1.23,
		B:       0x01,
		S:       1,
		I:       1,
		L:       1,
		D:       1.23,
	})

	doTestEncodeDecodeTuple(t, &Tuple{
		Byte:    math.MinInt8,
		Short:   math.MinInt16,
		Integer: math.MinInt32,
		Long:    math.MinInt64,
		Double:  -99.99,
		B:       0x00,
		S:       0,
		I:       0,
		L:       0,
		D:       -9999.9999,
	})

	doTestEncodeDecodeTuple(t, &Tuple{
		Byte:    math.MaxInt8,
		Short:   math.MaxInt16,
		Integer: math.MaxInt32,
		Long:    math.MaxInt64,
		Double:  math.MaxFloat32,
		B:       0xFF,
		S:       0xFFFF,
		I:       0xFFFFFFFF,
		L:       0xFFFFFFFFFFFFFFFF,
		D:       math.MaxFloat64,
	})
}

func doTestEncodeDecodeTuple(t *testing.T, tuple *Tuple) {
	e := NewEncoder()
	err := e.encObject(tuple)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	d := NewDecoder(e.buffer)
	decObj, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(tuple, decObj) {
		t.Errorf("expect: %v, but get: %v", tuple, decObj)
	}
}

type BasePointer struct {
	A *bool
}

func (t BasePointer) JavaClassName() string {
	return "test.base.Base"
}

func TestBasePointer(t *testing.T) {
	v := true
	base := BasePointer{
		A: &v,
	}
	doTestBasePointer(t, &base, &base)

	base = BasePointer{
		A: nil,
	}
	expectedF := false
	expectedBase := BasePointer{
		A: &expectedF,
	}
	doTestBasePointer(t, &base, &expectedBase)
}

func doTestBasePointer(t *testing.T, base *BasePointer, expected *BasePointer) {
	e := NewEncoder()
	err := e.encObject(base)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	d := NewDecoder(e.buffer)
	decObj, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(expected, decObj) {
		t.Errorf("expect: %v, but get: %v", base, decObj)
	}
}
