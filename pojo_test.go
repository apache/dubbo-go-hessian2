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
	"container/list"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_exception"
	"github.com/apache/dubbo-go-hessian2/java_util"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestGetGoName(t *testing.T) {
	assert.Equal(t, "time.Time", GetGoName(time.Now()))
	assert.Equal(t, "bytes.Buffer", GetGoName(bytes.Buffer{}))
	assert.Equal(t, "container/list/list.List", GetGoName(list.New()))
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/hessian.BusinessData", GetGoName(&BusinessData{}))
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/java_util/java_util.UUID", GetGoName(&java_util.UUID{}))
	assert.Equal(t, "github.com/apache/dubbo-go-hessian2/java_exception/java_exception.ClassNotFoundException", GetGoName(&java_exception.ClassNotFoundException{}))

	assert.Equal(t, "[]github.com/apache/dubbo-go-hessian2/hessian.BusinessData", GetGoName([]*BusinessData{}))
	assert.Equal(t, "[][]github.com/apache/dubbo-go-hessian2/hessian.BusinessData", GetGoName([][]*BusinessData{}))
}
