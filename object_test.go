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
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Department struct {
	Name string
}

// JavaClassName  java fully qualified path
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

// JavaClassName  java fully qualified path
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

// JavaClassName  java fully qualified path
func (UserName) JavaClassName() string {
	return "com.bdt.info.UserName"
}

type Person struct {
	UserName
	Age int32
	Sex bool
}

// JavaClassName  java fully qualified path
func (Person) JavaClassName() string {
	return "com.bdt.info.Person"
}

type JOB struct {
	Title   string
	Company string
}

// JavaClassName  java fully qualified path
func (JOB) JavaClassName() string {
	return "com.bdt.info.JOB"
}

type Worker struct {
	Person
	CurJob JOB
	Jobs   []JOB
}

// JavaClassName  java fully qualified path
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
			{Title: "manager", Company: "google"},
			{Title: "ceo", Company: "microsoft"},
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

func TestDecClassToMap(t *testing.T) {
	name := UserName{
		FirstName: "John",
		LastName:  "Doe",
	}
	person := Person{
		UserName: name,
		Age:      18,
		Sex:      true,
	}

	worker1 := &Worker{
		Person: person,
		CurJob: JOB{Title: "cto", Company: "facebook"},
		Jobs: []JOB{
			{Title: "manager", Company: "google"},
			{Title: "ceo", Company: "microsoft"},
		},
	}

	t.Logf("worker1: %v", worker1)

	e := NewEncoder()
	encErr := e.Encode(worker1)
	if encErr != nil {
		t.Fatalf("encode(worker:%#v) = error:%s", worker1, encErr)
	}
	data := e.Buffer()
	t.Logf("data: %s", data)

	// unRegisterPOJO before decode, so that to decode to map
	unRegisterPOJO(name)
	unRegisterPOJO(person)
	unRegisterPOJO(worker1)
	unRegisterPOJO(&worker1.Jobs[0])

	// strict mode
	d := NewDecoder(data)
	d.Strict = true
	res, err := d.Decode()
	if err == nil {
		t.Error("after unregister pojo, decoding should return error for strict mode")
		t.FailNow()
	}
	assert.Nil(t, res)

	// non-strict mode
	d = NewDecoder(data)
	res, err = d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("type of decode object:%v", reflect.TypeOf(res))

	worker2, ok := res.(map[string]interface{})
	if !ok {
		t.Fatalf("res:%#v should be a map for non-strict mode", res)
	}

	t.Logf("worker2: %v", worker2)

	// register pojo again
	RegisterPOJO(name)
	RegisterPOJO(person)
	RegisterPOJO(worker1)
	RegisterPOJO(&worker1.Jobs[0])

	// encode the map to object again
	e = NewEncoder()
	err = e.Encode(worker2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	data = e.Buffer()
	t.Logf("data: %s", data)

	// decode the encoded map data to struct
	d = NewDecoder(data)
	res, err = d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("type of decode object:%v", reflect.TypeOf(res))

	worker3, ok := res.(*Worker)
	if !ok {
		t.Fatalf("res:%#v should be a worker type", res)
	}

	t.Logf("worker3: %v", worker3)
	if !reflect.DeepEqual(worker1, worker3) {
		t.Fatal("worker1 not equal to worker3!")
	}
}

func TestEncodeMapToObject(t *testing.T) {
	name := &UserName{
		FirstName: "John",
		LastName:  "Doe",
	}

	// note: the first letter of the keys MUST lowercase.
	m := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	}

	e := NewEncoder()
	encErr := e.EncodeMapAsClass(name.JavaClassName(), m)
	if encErr != nil {
		t.Error(encErr)
		t.FailNow()
	}

	// register for decode map to object.
	RegisterPOJO(name)
	res := mustDecodeObject(t, e.Buffer())
	assert.True(t, reflect.DeepEqual(name, res))

	// unregister for encode map again.
	UnRegisterPOJOs(name)

	// note: the map contains the class key.
	m = map[string]interface{}{
		ClassKey:    name.JavaClassName(),
		"firstName": "John",
		"lastName":  "Doe",
	}

	// try to encode again
	e = NewEncoder()
	encErr = e.EncodeMapClass(m)
	if encErr != nil {
		t.Error(encErr)
		t.FailNow()
	}

	// register for decode map to object.
	RegisterPOJO(name)
	res = mustDecodeObject(t, e.Buffer())
	assert.True(t, reflect.DeepEqual(name, res))
}

type A0 struct{}

// JavaClassName  java fully qualified path
func (*A0) JavaClassName() string {
	return "com.caucho.hessian.test.A0"
}

type A1 struct{}

// JavaClassName  java fully qualified path
func (*A1) JavaClassName() string {
	return "com.caucho.hessian.test.A1"
}

type A2 struct{}

// JavaClassName  java fully qualified path
func (*A2) JavaClassName() string {
	return "com.caucho.hessian.test.A2"
}

type A3 struct{}

// JavaClassName  java fully qualified path
func (*A3) JavaClassName() string {
	return "com.caucho.hessian.test.A3"
}

type A4 struct{}

// JavaClassName  java fully qualified path
func (*A4) JavaClassName() string {
	return "com.caucho.hessian.test.A4"
}

type A5 struct{}

// JavaClassName  java fully qualified path
func (*A5) JavaClassName() string {
	return "com.caucho.hessian.test.A5"
}

type A6 struct{}

// JavaClassName  java fully qualified path
func (*A6) JavaClassName() string {
	return "com.caucho.hessian.test.A6"
}

type A7 struct{}

// JavaClassName  java fully qualified path
func (*A7) JavaClassName() string {
	return "com.caucho.hessian.test.A7"
}

type A8 struct{}

// JavaClassName  java fully qualified path
func (*A8) JavaClassName() string {
	return "com.caucho.hessian.test.A8"
}

type A9 struct{}

// JavaClassName  java fully qualified path
func (*A9) JavaClassName() string {
	return "com.caucho.hessian.test.A9"
}

type A10 struct{}

// JavaClassName  java fully qualified path
func (*A10) JavaClassName() string {
	return "com.caucho.hessian.test.A10"
}

type A11 struct{}

// JavaClassName  java fully qualified path
func (*A11) JavaClassName() string {
	return "com.caucho.hessian.test.A11"
}

type A12 struct{}

// JavaClassName  java fully qualified path
func (*A12) JavaClassName() string {
	return "com.caucho.hessian.test.A12"
}

type A13 struct{}

// JavaClassName  java fully qualified path
func (*A13) JavaClassName() string {
	return "com.caucho.hessian.test.A13"
}

type A14 struct{}

// JavaClassName  java fully qualified path
func (*A14) JavaClassName() string {
	return "com.caucho.hessian.test.A14"
}

type A15 struct{}

// JavaClassName  java fully qualified path
func (*A15) JavaClassName() string {
	return "com.caucho.hessian.test.A15"
}

type A16 struct{}

// JavaClassName  java fully qualified path
func (*A16) JavaClassName() string {
	return "com.caucho.hessian.test.A16"
}

type TestObjectStruct struct {
	Value int32 `hessian:"_value"`
}

// JavaClassName  java fully qualified path
func (*TestObjectStruct) JavaClassName() string {
	return "com.caucho.hessian.test.TestObject"
}

type TestConsStruct struct {
	First string          `hessian:"_first"`
	Rest  *TestConsStruct `hessian:"_rest"`
}

// JavaClassName  java fully qualified path
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

// JavaClassName  java fully qualified path
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

	testDecodeJavaData(t, "getTheTuple", "test.tuple.TupleProviderImpl", false, tuple)
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

// JavaClassName  java fully qualified path
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

func TestSkip(t *testing.T) {
	tempRegistry := pojoRegistry
	// recover the pojoRegistry after test, avoid influence other test cases.
	defer func() {
		pojoRegistry = tempRegistry
	}()

	// clear pojo
	pojoRegistry = &POJORegistry{
		j2g:      make(map[string]string),
		registry: make(map[string]*structInfo),
	}
	testDecodeFrameworkWithSkip(t, "replyObject_0", nil)
	testDecodeFrameworkWithSkip(t, "replyObject_1", nil)
	testDecodeFrameworkWithSkip(t, "replyObject_16", make([]interface{}, 17))
	testDecodeFrameworkWithSkip(t, "replyObject_2a", make([]interface{}, 2))
	testDecodeFrameworkWithSkip(t, "replyObject_3", nil)

	testDecodeFrameworkWithSkip(t, "replyTypedMap_0", make(map[interface{}]interface{}))

	testDecodeFrameworkWithSkip(t, "replyTypedFixedList_0", make([]string, 0))
	testDecodeFrameworkWithSkip(t, "replyUntypedFixedList_0", []interface{}{})

	testDecodeFrameworkWithSkip(t, "customReplyTypedFixedListHasNull", make([]Object, 3))
	testDecodeFrameworkWithSkip(t, "customReplyTypedVariableListHasNull", make([]Object, 3))
	testDecodeFrameworkWithSkip(t, "customReplyUntypedFixedListHasNull", make([]interface{}, 3))
	testDecodeFrameworkWithSkip(t, "customReplyUntypedVariableListHasNull", make([]interface{}, 3))

	testDecodeFrameworkWithSkip(t, "customReplyTypedFixedList_A0", make([]interface{}, 3))
	testDecodeFrameworkWithSkip(t, "customReplyTypedVariableList_A0", make([]interface{}, 3))

	testDecodeFrameworkWithSkip(t, "customReplyTypedFixedList_Test", nil)

	testDecodeFrameworkWithSkip(t, "customReplyTypedFixedList_Object", make([]Object, 1))
}

type Animal struct {
	Name string
}

type animal struct {
	Name string
}

func (a Animal) JavaClassName() string {
	return "test.Animal"
}

type Dog struct {
	Animal
	animal
	Gender  string
	DogName string `hessian:"-"`
}

// JavaClassName  java fully qualified path
func (dog Dog) JavaClassName() string {
	return "test.Dog"
}

type DogAll struct {
	All    bool
	Name   string
	Gender string
}

// JavaClassName  java fully qualified path
func (dog *DogAll) JavaClassName() string {
	return "test.DogAll"
}

// see https://github.com/apache/dubbo-go-hessian2/issues/149
func TestIssue149_EmbedStructGoDecode(t *testing.T) {
	t.Run(`extends to embed`, func(t *testing.T) {
		RegisterPOJO(&Dog{})
		got, err := decodeJavaResponse(`customReplyExtendClass`, ``, false)
		if err != nil {
			t.Error(err)
		}

		want := &Dog{Animal{`a dog`}, animal{}, `male`, ``}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v got %v", want, got)
		}
	})

	t.Run(`extends to all fields`, func(t *testing.T) {
		RegisterPOJO(&DogAll{})
		got, err := decodeJavaResponse(`customReplyExtendClassToSingleStruct`, ``, false)
		if err != nil {
			t.Error(err)
		}

		want := &DogAll{true, `a dog`, `male`}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v got %v", want, got)
		}
	})
}

func TestIssue150_EmbedStructJavaDecode(t *testing.T) {
	RegisterPOJO(&Dog{})
	RegisterPOJO(&Animal{})

	dog := &Dog{Animal{`a dog`}, animal{}, `male`, `DogName`}
	bytes, err := encodeTarget(dog)
	t.Log(string(bytes), err)

	testJavaDecode(t, "customArgTypedFixed_Extends", dog)
}

type Mix struct {
	A  int
	B  string
	CA time.Time
	CB int64
	CC string
	CD []float64
	D  map[string]interface{}
}

func (m Mix) JavaClassName() string {
	return `test.mix`
}

func init() {
	RegisterPOJO(new(Mix))
}

// BenchmarkJsonEncode-8   	  217354	      4799 ns/op	     832 B/op	      15 allocs/op
func BenchmarkJsonEncode(b *testing.B) {
	m := Mix{A: int('a'), B: `hello`}
	m.CD = []float64{1, 2, 3}
	m.D = map[string]interface{}{`floats`: m.CD, `A`: m.A, `m`: m}

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&m)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkEncode-8   	  211452	      5560 ns/op	    1771 B/op	      51 allocs/op
func BenchmarkEncode(b *testing.B) {
	m := Mix{A: int('a'), B: `hello`}
	m.CD = []float64{1, 2, 3}
	m.D = map[string]interface{}{`floats`: m.CD, `A`: m.A, `m`: m}

	for i := 0; i < b.N; i++ {
		_, err := encodeTarget(&m)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkJsonDecode-8   	  123922	      8549 ns/op	    1776 B/op	      51 allocs/op
func BenchmarkJsonDecode(b *testing.B) {
	m := Mix{A: int('a'), B: `hello`}
	m.CD = []float64{1, 2, 3}
	m.D = map[string]interface{}{`floats`: m.CD, `A`: m.A, `m`: m}
	bytes, err := json.Marshal(&m)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		m := &Mix{}
		err := json.Unmarshal(bytes, m)
		if err != nil {
			b.Error(err)
		}
	}
}

// BenchmarkDecode-8   	  104196	     10924 ns/op	    6424 B/op	      98 allocs/op
func BenchmarkDecode(b *testing.B) {
	m := Mix{A: int('a'), B: `hello`}
	m.CD = []float64{1, 2, 3}
	m.D = map[string]interface{}{`floats`: m.CD, `A`: m.A, `m`: m}
	bytes, err := encodeTarget(&m)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		d := NewDecoder(bytes)
		_, err := d.Decode()
		if err != nil {
			b.Error(err)
		}
	}
}

type Person183 struct {
	Name string
}

func (Person183) JavaClassName() string {
	return `test.Person183`
}

func TestIssue183_DecodeExcessStructField(t *testing.T) {
	RegisterPOJO(&Person183{})
	got, err := decodeJavaResponse(`customReplyPerson183`, ``, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
}

type GenericResponse struct {
	Code int
	Data interface{}
}

func (GenericResponse) JavaClassName() string {
	return `test.generic.Response`
}

type BusinessData struct {
	Name  string
	Count int
}

func (BusinessData) JavaClassName() string {
	return `test.generic.BusinessData`
}

func TestCustomReplyGenericResponseLong(t *testing.T) {
	res := &GenericResponse{
		Code: 200,
		Data: int64(123),
	}
	RegisterPOJO(res)

	testDecodeFramework(t, "customReplyGenericResponseLong", res)
}

func TestCustomReplyGenericResponseBusinessData(t *testing.T) {
	data := BusinessData{
		Name:  "apple",
		Count: 5,
	}
	res := &GenericResponse{
		Code: 201,
		Data: data,
	}
	RegisterPOJO(data)
	RegisterPOJO(res)

	testDecodeFramework(t, "customReplyGenericResponseBusinessData", res)
}

func TestCustomReplyGenericResponseList(t *testing.T) {
	data := []*BusinessData{
		{
			Name:  "apple",
			Count: 5,
		},
		{
			Name:  "banana",
			Count: 6,
		},
	}
	res := &GenericResponse{
		Code: 202,
		Data: data,
	}
	RegisterPOJO(data[0])
	RegisterPOJO(res)

	testDecodeFrameworkFunc(t, "customReplyGenericResponseList", func(r interface{}) {
		expect, ok := r.(*GenericResponse)
		if !ok {
			t.Errorf("expect *GenericResponse, but get %v", r)
			return
		}
		list, dataOk := expect.Data.([]interface{})
		if !dataOk {
			t.Errorf("expect []interface{}, but get %v", expect.Data)
			return
		}
		assert.Equal(t, res.Code, expect.Code)
		assert.True(t, reflect.DeepEqual(data[0], list[0]))
		assert.True(t, reflect.DeepEqual(data[1], list[1]))
	})
}

func TestWrapperClassArray(t *testing.T) {
	t.Log(getListType("java.lang.Byte"))
	t.Log(getListType("java.lang.Integer"))
	got, err := decodeJavaResponse(`byteArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var b1, b2, b3 byte = 'A', 'B', 'C'
	ba := []*byte{&b1, &b2, &b3}
	assert.True(t, reflect.DeepEqual(got, ba))

	got, err = decodeJavaResponse(`shortArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var s1, s2, s3 int16 = 1, 100, 10000
	sa := []*int16{&s1, &s2, &s3}
	assert.True(t, reflect.DeepEqual(got, sa))

	got, err = decodeJavaResponse(`integerArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var i1, i2, i3 int32 = 1, 100, 10000
	ia := []*int32{&i1, &i2, &i3}
	assert.True(t, reflect.DeepEqual(got, ia))

	got, err = decodeJavaResponse(`longArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var l1, l2, l3 int64 = 1, 100, 10000
	la := []*int64{&l1, &l2, &l3}
	assert.True(t, reflect.DeepEqual(got, la))

	got, err = decodeJavaResponse(`characterArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11 Rune = 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'
	ca := []*int32{(*int32)(&r1), (*int32)(&r2), (*int32)(&r3), (*int32)(&r4), (*int32)(&r5), (*int32)(&r6), (*int32)(&r7), (*int32)(&r8), (*int32)(&r9), (*int32)(&r10), (*int32)(&r11)}
	assert.True(t, reflect.DeepEqual(got, ca))

	got, err = decodeJavaResponse(`booleanArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var bl1, bl2, bl3 bool = true, false, true
	bla := []*bool{&bl1, &bl2, &bl3}
	assert.True(t, reflect.DeepEqual(got, bla))

	got, err = decodeJavaResponse(`floatArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var f1, f2, f3 float32 = 1.0, 100.0, 10000.1
	fa := []*float32{&f1, &f2, &f3}
	assert.True(t, reflect.DeepEqual(got, fa))

	got, err = decodeJavaResponse(`doubleArray`, `test.TestWrapperClassArray`, false)
	assert.NoError(t, err)
	t.Logf("%T %+v", got, got)
	var fl1, fl2, fl3 float64 = 1.0, 100.0, 10000.1
	da := []*float64{&fl1, &fl2, &fl3}
	assert.True(t, reflect.DeepEqual(got, da))
}

type User struct {
	Id   *int32
	List []int32
}

func (u *User) JavaClassName() string {
	return "test.model.User"
}

func TestDecodeIntegerHasNull(t *testing.T) {
	RegisterPOJO(&User{})
	id := int32(1)
	testDecodeFramework(t, "customReplyTypedIntegerHasNull", &User{Id: &id})
}

func TestDecodeSliceIntegerHasNull(t *testing.T) {
	RegisterPOJO(&User{})
	id := int32(1)
	testDecodeFramework(t, "customReplyTypedListIntegerHasNull", &User{Id: &id, List: []int32{1, 0}})
}

func TestDecodeCustomReplyEnumVariableList(t *testing.T) {
	for v := range _LocaleCategoryEnumValues {
		RegisterJavaEnum(v)
	}

	got, err := decodeJavaResponse(`customReplyEnumVariableList`, ``, false)
	assert.NoError(t, err)
	t.Logf("customReplyEnumVariableList %T %+v", got, got)
	enumList, ok := got.([]*LocaleCategoryEnum)
	if !ok {
		t.Errorf("expect []LocaleCategoryEnum, but get %v", got)
		return
	}
	assert.Equal(t, LocaleCategoryEnumDisplay, *enumList[0])
	assert.Nil(t, enumList[1])
	assert.Equal(t, LocaleCategoryEnumFormat, *enumList[2])
}

const (
	LocaleCategoryEnumDisplay LocaleCategoryEnum = iota
	LocaleCategoryEnumFormat
)

var (
	_LocaleCategoryEnumValues = map[LocaleCategoryEnum]string{
		LocaleCategoryEnumDisplay: "DISPLAY",
		LocaleCategoryEnumFormat:  "FORMAT",
	}
	_LocaleCategoryEnumEntities = map[string]LocaleCategoryEnum{
		"DISPLAY": LocaleCategoryEnumDisplay,
		"FORMAT":  LocaleCategoryEnumFormat,
	}
)

type LocaleCategoryEnum JavaEnum

func (e LocaleCategoryEnum) JavaClassName() string {
	return "java.util.Locale$Category"
}

func (e LocaleCategoryEnum) String() string {
	if v, ok := _LocaleCategoryEnumValues[e]; ok {
		return v
	}

	return strconv.Itoa(int(e))
}

func (e LocaleCategoryEnum) EnumValue(s string) JavaEnum {
	if v, ok := _LocaleCategoryEnumEntities[s]; ok {
		return JavaEnum(v)
	}

	return InvalidJavaEnum
}
