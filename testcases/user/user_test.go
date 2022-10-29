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

package user

import (
	"testing"
	"time"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/stretchr/testify/assert"
)

func TestEnumConvert(t *testing.T) {
	var g interface{}
	g = WOMAN

	// new defined type cant be converted to the original type.
	failConvertedValue, ok := g.(hessian.JavaEnum)
	assert.False(t, ok)
	assert.Equal(t, hessian.JavaEnum(0), failConvertedValue)
}

func TestUserEncodeDecode(t *testing.T) {
	ts, _ := time.Parse("2006-01-02 15:04:05", "2019-01-01 12:34:56")
	u1 := &User{ID: "001", Name: "Lily", Age: 18, Time: ts.Local(), Sex: WOMAN}
	hessian.RegisterPOJO(u1)

	encoder := hessian.NewEncoder()
	err := encoder.Encode(u1)
	assert.Nil(t, err)

	buf := encoder.Buffer()
	decoder := hessian.NewDecoder(buf)
	dec, err := decoder.Decode()
	assert.Nil(t, err)
	assert.Equal(t, u1, dec)
}
