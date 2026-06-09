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

package java_exception

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDubboGenericExceptionErrorMessage(t *testing.T) {
	tests := []struct {
		name             string
		exceptionClass   string
		exceptionMessage string
		wantError        string
	}{
		{
			name:             "class and message",
			exceptionClass:   "com.example.UserNotFoundException",
			exceptionMessage: "user not found",
			wantError:        "java exception: com.example.UserNotFoundException - user not found",
		},
		{
			name:             "message only",
			exceptionMessage: "user not found",
			wantError:        "user not found",
		},
		{
			name:           "class only",
			exceptionClass: "com.example.UserNotFoundException",
			wantError:      "com.example.UserNotFoundException",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exception := NewDubboGenericException(test.exceptionClass, test.exceptionMessage)
			assert.Equal(t, test.exceptionClass, exception.ExceptionClass)
			assert.Equal(t, test.exceptionMessage, exception.ExceptionMessage)
			assert.Equal(t, test.wantError, exception.DetailMessage)
			assert.Equal(t, test.wantError, exception.Error())
		})
	}
}

func TestDubboGenericExceptionErrorFallback(t *testing.T) {
	exception := DubboGenericException{
		ExceptionClass:   "com.example.UserNotFoundException",
		ExceptionMessage: "user not found",
	}

	assert.Equal(t, "java exception: com.example.UserNotFoundException - user not found", exception.Error())
}

func TestDubboGenericExceptionErrorPrefersDetailMessage(t *testing.T) {
	exception := DubboGenericException{
		DetailMessage:    "decoded detail",
		ExceptionClass:   "com.example.UserNotFoundException",
		ExceptionMessage: "user not found",
	}

	assert.Equal(t, "decoded detail", exception.Error())
}
