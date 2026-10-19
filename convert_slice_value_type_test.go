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
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestConvertSliceValueTypeFloat(t *testing.T) {
	v, err := ConvertSliceValueType(reflect.TypeOf([]float64{}), reflect.ValueOf([]interface{}{1.5, 2.5}))
	assert.Nil(t, err)
	assert.Equal(t, []float64{1.5, 2.5}, v.Interface())
}

func TestConvertSliceValueTypeFloatKindMismatch(t *testing.T) {
	_, err := ConvertSliceValueType(reflect.TypeOf([]float64{}), reflect.ValueOf([]interface{}{"abc"}))
	assert.NotNil(t, err)
}

func TestConvertSliceValueTypeInt(t *testing.T) {
	v, err := ConvertSliceValueType(reflect.TypeOf([]int64{}), reflect.ValueOf([]interface{}{int64(1), int64(2)}))
	assert.Nil(t, err)
	assert.Equal(t, []int64{1, 2}, v.Interface())
}

func TestConvertSliceValueTypeIntKindMismatch(t *testing.T) {
	_, err := ConvertSliceValueType(reflect.TypeOf([]int64{}), reflect.ValueOf([]interface{}{"abc"}))
	assert.NotNil(t, err)
}

func TestConvertSliceValueTypeUintFromInt(t *testing.T) {
	v, err := ConvertSliceValueType(reflect.TypeOf([]uint64{}), reflect.ValueOf([]interface{}{int64(1), int64(2)}))
	assert.Nil(t, err)
	assert.Equal(t, []uint64{1, 2}, v.Interface())
}

func TestConvertSliceValueTypeUintKindMismatch(t *testing.T) {
	_, err := ConvertSliceValueType(reflect.TypeOf([]uint64{}), reflect.ValueOf([]interface{}{"abc"}))
	assert.NotNil(t, err)
}
