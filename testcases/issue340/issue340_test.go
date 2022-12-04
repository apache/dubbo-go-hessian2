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

package issue340

import (
	"reflect"
	"testing"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/stretchr/testify/assert"
)

type Point struct {
	X int
	Y int
}

func (Point) JavaClassName() string {
	return "com.test.Point"
}

type SpPoint struct {
	X  int
	Y  int
	Sp int
}

func (SpPoint) JavaClassName() string {
	return "com.test.SpPoint"
}

type ReqInfo struct {
	Name     string
	Points   []*Point
	SpPoints []*SpPoint
}

func (ReqInfo) JavaClassName() string {
	return "com.test.ReqInfo"
}

func TestIssue340Case(t *testing.T) {
	req := &ReqInfo{
		Name:     "test",
		Points:   []*Point{},
		SpPoints: []*SpPoint{},
	}

	hessian.RegisterPOJO(&Point{})
	hessian.RegisterPOJO(&SpPoint{})
	hessian.RegisterPOJO(req)

	encoder := hessian.NewEncoder()
	err := encoder.Encode(req)
	if err != nil {
		t.Error(err)
		return
	}

	enBuf := encoder.Buffer()

	decoder := hessian.NewDecoder(enBuf)
	deReq, err := decoder.Decode()
	assert.Nil(t, err)

	t.Log(deReq)

	assert.True(t, reflect.DeepEqual(req, deReq))
}
