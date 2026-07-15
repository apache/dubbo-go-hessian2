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
	"strings"
)

import (
	"github.com/apache/dubbo-go-hessian2/java_exception"
)

// GenericException represents a Java exception with its class name and message.
//
// Previously, when dubbo-go invoked a Java service and the Java side threw
// a business exception (e.g. UserNotFoundException), the original exception
// type was lost — the Go caller only received a plain string error message
// and could not determine which specific exception was thrown.
//
// This type preserves the full exception information:
//   - ExceptionClass:  the fully qualified Java class name, e.g. "com.example.UserNotFoundException"
//   - ExceptionMessage: the exception detail message, e.g. "user not found"
//
// See also:
//   - Java dubbo's GenericException.java (dubbo-common/src/main/java/org/apache/dubbo/rpc/service/GenericException.java)
//   - GitHub issue: https://github.com/apache/dubbo-go/issues/3167
type GenericException struct {
	ExceptionClass   string
	ExceptionMessage string
}

// Error returns a readable error string.
func (e GenericException) Error() string {
	if e.ExceptionClass == "" {
		return e.ExceptionMessage
	}
	if e.ExceptionMessage == "" {
		return e.ExceptionClass
	}
	return "java exception: " + e.ExceptionClass + " - " + e.ExceptionMessage
}

// ToGenericException converts decoded exception to GenericException when possible.
func ToGenericException(expt any) (*GenericException, bool) {
	switch v := expt.(type) {
	case *GenericException:
		if v == nil {
			return nil, false
		}
		return v, true
	case GenericException:
		return &v, true
	case *java_exception.DubboGenericException:
		if v == nil {
			return nil, false
		}
		return &GenericException{ExceptionClass: v.ExceptionClass, ExceptionMessage: v.ExceptionMessage}, true
	case java_exception.DubboGenericException:
		return &GenericException{ExceptionClass: v.ExceptionClass, ExceptionMessage: v.ExceptionMessage}, true
	case java_exception.Throwabler:
		if rv := reflect.ValueOf(v); rv.Kind() == reflect.Ptr && rv.IsNil() {
			return nil, false
		}
		return &GenericException{ExceptionClass: v.JavaClassName(), ExceptionMessage: v.Error()}, true
	case string:
		return parseLegacyException(v), true
	}
	return nil, false
}

func parseLegacyException(exStr string) *GenericException {
	const prefix = "java exception:"
	msg := strings.TrimSpace(exStr)
	if strings.HasPrefix(msg, prefix) {
		msg = strings.TrimSpace(strings.TrimPrefix(msg, prefix))
		if class, message, ok := strings.Cut(msg, " - "); ok {
			return &GenericException{
				ExceptionClass:   strings.TrimSpace(class),
				ExceptionMessage: strings.TrimSpace(message),
			}
		}
	}
	return &GenericException{ExceptionClass: "java.lang.Exception", ExceptionMessage: msg}
}
