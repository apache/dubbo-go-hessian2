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
	"testing"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_exception"
	"github.com/stretchr/testify/assert"
)

func TestParseLegacyException(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantClass   string
		wantMessage string
	}{
		{
			name:        "plain string",
			input:       "user not found",
			wantClass:   "java.lang.Exception",
			wantMessage: "user not found",
		},
		{
			name:        "legacy prefix without separator",
			input:       "java exception: user not found",
			wantClass:   "java.lang.Exception",
			wantMessage: "user not found",
		},
		{
			name:        "full Error() format",
			input:       "java exception: com.example.FooException - something went wrong",
			wantClass:   "com.example.FooException",
			wantMessage: "something went wrong",
		},
		{
			name:        "message with separator inside",
			input:       "java exception: com.example.FooException - error - more details",
			wantClass:   "com.example.FooException",
			wantMessage: "error - more details",
		},
		{
			name:        "prefix with trailing spaces",
			input:       "  java exception:  com.example.FooException  -  something went wrong  ",
			wantClass:   "com.example.FooException",
			wantMessage: "something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ge := parseLegacyException(tt.input)
			assert.Equal(t, tt.wantClass, ge.ExceptionClass)
			assert.Equal(t, tt.wantMessage, ge.ExceptionMessage)
		})
	}
}

func TestGenericExceptionError(t *testing.T) {
	tests := []struct {
		name      string
		class     string
		message   string
		wantError string
	}{
		{
			name:      "class and message",
			class:     "com.example.FooException",
			message:   "something went wrong",
			wantError: "java exception: com.example.FooException - something went wrong",
		},
		{
			name:      "message only",
			message:   "user not found",
			wantError: "user not found",
		},
		{
			name:      "class only",
			class:     "com.example.FooException",
			wantError: "com.example.FooException",
		},
		{
			name:      "empty both",
			wantError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := GenericException{ExceptionClass: tt.class, ExceptionMessage: tt.message}
			assert.Equal(t, tt.wantError, e.Error())
		})
	}
}

func TestToGenericException(t *testing.T) {
	t.Run("from *GenericException", func(t *testing.T) {
		ge := &GenericException{ExceptionClass: "com.example.Foo", ExceptionMessage: "bar"}
		got, ok := ToGenericException(ge)
		assert.True(t, ok)
		assert.Same(t, ge, got)
	})

	t.Run("from GenericException value", func(t *testing.T) {
		ge := GenericException{ExceptionClass: "com.example.Foo", ExceptionMessage: "bar"}
		got, ok := ToGenericException(ge)
		assert.True(t, ok)
		assert.Equal(t, ge.ExceptionClass, got.ExceptionClass)
		assert.Equal(t, ge.ExceptionMessage, got.ExceptionMessage)
	})

	t.Run("from *DubboGenericException", func(t *testing.T) {
		dge := &java_exception.DubboGenericException{
			ExceptionClass:   "com.example.Foo",
			ExceptionMessage: "bar",
		}
		got, ok := ToGenericException(dge)
		assert.True(t, ok)
		assert.Equal(t, "com.example.Foo", got.ExceptionClass)
		assert.Equal(t, "bar", got.ExceptionMessage)
	})

	t.Run("from DubboGenericException value", func(t *testing.T) {
		dge := java_exception.DubboGenericException{
			ExceptionClass:   "com.example.Foo",
			ExceptionMessage: "bar",
		}
		got, ok := ToGenericException(dge)
		assert.True(t, ok)
		assert.Equal(t, "com.example.Foo", got.ExceptionClass)
		assert.Equal(t, "bar", got.ExceptionMessage)
	})

	t.Run("from Throwabler", func(t *testing.T) {
		thr := java_exception.NewThrowable("some error")
		got, ok := ToGenericException(thr)
		assert.True(t, ok)
		assert.Equal(t, thr.JavaClassName(), got.ExceptionClass)
		assert.Equal(t, thr.Error(), got.ExceptionMessage)
	})

	t.Run("from string", func(t *testing.T) {
		got, ok := ToGenericException("user not found")
		assert.True(t, ok)
		assert.Equal(t, "java.lang.Exception", got.ExceptionClass)
		assert.Equal(t, "user not found", got.ExceptionMessage)
	})

	t.Run("from string with Error() format", func(t *testing.T) {
		got, ok := ToGenericException("java exception: com.example.Foo - something went wrong")
		assert.True(t, ok)
		assert.Equal(t, "com.example.Foo", got.ExceptionClass)
		assert.Equal(t, "something went wrong", got.ExceptionMessage)
	})

	t.Run("from unsupported type", func(t *testing.T) {
		got, ok := ToGenericException(42)
		assert.False(t, ok)
		assert.Nil(t, got)
	})
}
