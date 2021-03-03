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

func init() {
	SetCollectionSerialize(&JavaHashSet{})
}

type JavaHashSet struct {
	value []interface{}
}

func (j *JavaHashSet) Get() []interface{} {
	return j.value
}

func (j *JavaHashSet) Set(v []interface{}) {
	j.value = v
}

func (j *JavaHashSet) JavaClassName() string {
	return "java.util.HashSet"
}

func TestListJavaCollectionEncode(t *testing.T) {
	inside := make([]interface{}, 2)
	inside[0] = int32(0)
	inside[1] = int32(1)
	hashSet := JavaHashSet{value: inside}
	testJavaDecode(t, "customArgTypedFixedList_HashSet", &hashSet)
}

func TestListJavaCollectionDecode(t *testing.T) {
	inside := make([]interface{}, 2)
	inside[0] = int32(0)
	inside[1] = int32(1)
	hashSet := JavaHashSet{value: inside}
	testDecodeFramework(t, "customReplyTypedFixedList_HashSet", &hashSet)
}
