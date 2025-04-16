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

// Unit test for decoding hessian2 based on official api with doc on
// http://javadoc4.caucho.com/com/caucho/hessian/test/TestHessian2.html.
// One can call the api by running the local test_hessian.jar or sending
// a request to the remote server http://hessian.caucho.com/test/test
// directly.

package hessian

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"unsafe"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_exception"
)

const (
	hessianJar = "test_hessian/target/test_hessian-1.0.0.jar"
	testString = "hello, world! 你好，世界！"
)

func isFileExist(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func genHessianJar() {
	existFlag := isFileExist(hessianJar)
	if existFlag {
		return
	}

	log.Println("generate hessian jar")
	cmd := exec.Command("mvn", "clean", "package")
	cmd.Dir = "./test_hessian"
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("after exec command 'mvn clean package', got error:%v, error output:%v",
			err, string(out))
	}
}

func getJavaReply(method, className string) []byte {
	genHessianJar()
	log.Println("get java reply: ", className, method)
	cmdArgs := []string{"-jar", hessianJar, method}
	if className != "" {
		cmdArgs = append(cmdArgs, className)
	}
	cmd := exec.Command("java", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Java command exited with status %d: %s", exitErr.ExitCode(), string(out))
			return out
		}
		log.Fatal(cmd.Args, err)
	}
	return out
}

func decodeJavaResponse(method, className string, skip bool) (interface{}, error) {
	b := getJavaReply(method, className)
	var d *Decoder
	if skip {
		d = NewDecoderWithSkip(b)
	} else {
		d = NewDecoder(b)
	}
	r, e := d.Decode()
	if e != nil {
		return nil, e
	}
	return r, nil
}

func testDecodeFramework(t *testing.T, method string, expected interface{}) {
	testDecodeJavaData(t, method, "", false, expected)
}

func testDecodeFrameworkWithSkip(t *testing.T, method string, expected interface{}) {
	testDecodeJavaData(t, method, "", true, expected)
}

func testDecodeJavaData(t *testing.T, method, className string, skip bool, expected interface{}) {
	r, e := decodeJavaResponse(method, className, skip)
	if e != nil {
		t.Errorf("%s: decode fail with error: %v", method, e)
		return
	}

	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	trow, o1 := r.(java_exception.Throwabler)
	expe, o2 := expected.(java_exception.Throwabler)
	if o1 && o2 {
		log.Println(reflect.TypeOf(trow), reflect.TypeOf(trow).Elem().Name())
		if trow.Error() == expe.Error() && reflect.TypeOf(trow).Elem().Name() == reflect.TypeOf(expe).Elem().Name() {
			return
		}
		t.Errorf("%s: got %v, wanted %v", method, r, expected)
	} else if !reflect.DeepEqual(r, expected) {
		t.Errorf("%s: got %v, wanted %v", method, r, expected)
	}
}

func testDecodeFrameworkFunc(t *testing.T, method string, expected func(interface{})) {
	r, e := decodeJavaResponse(method, "", false)
	if e != nil {
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	expected(r)
}

func mustDecodeObject(t *testing.T, b []byte) interface{} {
	d := NewDecoder(b)
	res, err := d.Decode()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return res
}

func TestUserDefinedException(t *testing.T) {
	expect := &UnknownException{
		DetailMessage: "throw UserDefinedException",
	}
	testDecodeFramework(t, "throw_UserDefinedException", expect)
}

type Circular214 struct {
	Num      int
	Previous *Circular214
	Next     *Circular214
	Bytes    []byte
}

func (Circular214) JavaClassName() string {
	return "com.company.Circular"
}

func (c *Circular214) String() string {
	return fmt.Sprintf("Addr:%p, Num: %d, Previous: %p, Next: %p, Bytes: %s", c, c.Num, c.Previous, c.Next, c.Bytes)
}

func TestIssue214(t *testing.T) {
	c := &Circular214{}
	c.Num = 1234
	c.Previous = c
	c.Next = c
	c.Bytes = []byte(`{"a":"b"}`)
	e := NewEncoder()
	err := e.Encode(c)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
		return
	}

	bytes := e.Buffer()
	decoder := NewDecoder(bytes)
	decode, err := decoder.Decode()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
		return
	}
	t.Log(decode)
	assert.True(t, reflect.DeepEqual(c, decode))
}

type Issue299Args1 struct {
	Label string
	Key   string
}

func (Issue299Args1) JavaClassName() string {
	return "com.test.Issue299Args1"
}

type Issue299MockData struct {
	SliceArgs []interface{}
	MapArgs   map[string]interface{}
}

func (Issue299MockData) JavaClassName() string {
	return "com.test.Issue299MockData"
}

func TestIssue299HessianDecode(t *testing.T) {
	RegisterPOJO(new(Issue299Args1))
	RegisterPOJO(new(Issue299MockData))

	d := &Issue299MockData{
		SliceArgs: []interface{}{
			[]*Issue299Args1{
				{Label: "1", Key: "2"},
			},
		},
		MapArgs: map[string]interface{}{
			"interface_slice": []interface{}{
				"MyName",
			},
			"interface_map": map[interface{}]interface{}{
				"k1": "v1",
			},
		},
	}

	encoder := NewEncoder()
	err := encoder.Encode(d)
	if err != nil {
		t.Errorf("encode obj error: %v", err)
		return
	}
	decoder := NewDecoder(encoder.Buffer())
	doInterface, err := decoder.Decode()
	if err != nil {
		t.Errorf("decode obj error: %v", err)
		return
	}
	do := doInterface.(*Issue299MockData)
	if !reflect.DeepEqual(d, do) {
		t.Errorf("not equal d: %#v, do: %#v", d, do)
		return
	}
}

type Issue323B struct {
	Num int
}

func (b *Issue323B) JavaClassName() string {
	return "B"
}

type Issue323BB struct {
	List1 []*Issue323B
	List2 []*Issue323B
}

func (bb *Issue323BB) JavaClassName() string {
	return "BB"
}

func TestIssue323(t *testing.T) {
	RegisterPOJO(&Issue323B{})
	RegisterPOJO(&Issue323BB{})
	a1 := &Issue323B{
		Num: 1,
	}
	a2 := &Issue323B{
		Num: 2,
	}
	list := []*Issue323B{a1, a2}
	b := &Issue323BB{
		List1: list,
		List2: list,
	}
	e := NewEncoder()
	err := e.Encode(b)
	assert.Nil(t, err)
	fmt.Println(b)

	d := NewDecoder(e.Buffer())
	res, err := d.Decode()
	assert.Nil(t, err)
	fmt.Println(res)
	assert.True(t, reflect.DeepEqual(b, res))

	decodB, ok := res.(*Issue323BB)
	if !ok {
		t.Log("res is not Issue323BB")
		t.FailNow()
	}

	// list1 and list2 should be reference to the same one.
	assert.Equal(t, unsafe.Pointer(reflect.ValueOf(decodB.List1).Pointer()),
		unsafe.Pointer(reflect.ValueOf(decodB.List2).Pointer()))
}
