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

// Unit test for decoding hessian2 based on official api with doc on
// http://javadoc4.caucho.com/com/caucho/hessian/test/TestHessian2.html.
// One can call the api by running the local test_hessian.jar or sending
// a request to the remote server http://hessian.caucho.com/test/test
// directly.
package hessian

import (
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

const (
	hessianJar = "test_hessian/target/test_hessian-1.0.0.jar"
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
		log.Fatal(err)
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
		t.Errorf("%s: decode fail with error %v", method, e)
		return
	}

	tmp, ok := r.(*_refHolder)
	if ok {
		r = tmp.value.Interface()
	}
	if !reflect.DeepEqual(r, expected) {
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
