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

package hessian_test

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-hessian2"
	dupclass "github.com/apache/dubbo-go-hessian2/hessian_test/hessian_test"
)

import (
	"github.com/stretchr/testify/assert"
)

const (
	ExpectedErrorMsg = "reflect.Set: value of type hessian_test.CaseZ is not assignable to type hessian_test.CaseZ"
)

type CaseZ struct {
	Name string
}

func (CaseZ) JavaClassName() string {
	return "com.test.caseZ"
}

func TestDuplicatedClassGetGoType(t *testing.T) {
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/hessian_test_test/hessian_test.CaseZ", hessian.GetGoType(&CaseZ{}))
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/hessian_test/hessian_test/hessian_test.CaseZ", hessian.GetGoType(&dupclass.CaseZ{}))
}
