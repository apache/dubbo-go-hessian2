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
	"bytes"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

var assertEqual = func(want, got []byte, t *testing.T) {
	if !bytes.Equal(want, got) {
		t.Fatalf("want %v , got %v", want, got)
	}
}

func encodeTarget(target interface{}) ([]byte, error) {
	e := NewEncoder()
	err := e.Encode(target)
	if err != nil {
		return nil, err
	}
	return e.Buffer(), nil
}

func javaDecodeValidate(t *testing.T, method string, target interface{}) (string, error) {
	b, err := encodeTarget(target)
	if err != nil {
		return "", err
	}

	genHessianJar()
	cmd := exec.Command("java", "-jar", hessianJar, method)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Logf("call java error: %v", err)
		return "", err
	}

	go func() {
		_, _ = stdin.Write(b)
		_ = stdin.Close()
	}()

	out, err := cmd.Output()
	if err != nil {
		t.Logf("get java result error: %v", err)
		return "", err
	}
	return string(out), nil
}

func testJavaDecode(t *testing.T, method string, target interface{}) {
	result, err := javaDecodeValidate(t, method, target)
	if err != nil {
		t.Errorf("%s: encode fail with error: %v", method, err)
		return
	}

	if result != "true" {
		t.Errorf("%s: encode %v to bytes wrongly, result: %s", method, target, result)
	}
}

func testSimpleEncode(t *testing.T, v interface{}) {
	e := NewEncoder()
	err := e.Encode(v)
	assert.Nil(t, err)
}
