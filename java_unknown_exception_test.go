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
	"github.com/stretchr/testify/assert"
)

func TestCheckAndGetException(t *testing.T) {
	clazzInfo1 := classInfo{
		javaName:      "com.test.UserDefinedException",
		fieldNameList: []string{"detailMessage", "code", "suppressedExceptions", "stackTrace", "cause"},
	}
	s, b := checkAndGetException(clazzInfo1)
	assert.True(t, b)

	assert.Equal(t, s.javaName, "com.test.UserDefinedException")
	assert.Equal(t, s.goName, "hessian.UnknownException")

	clazzInfo2 := classInfo{
		javaName:      "com.test.UserDefinedException",
		fieldNameList: []string{"detailMessage", "code", "suppressedExceptions", "cause"},
	}
	s, b = checkAndGetException(clazzInfo2)
	assert.False(t, b)
	assert.Equal(t, s, structInfo{})
}
