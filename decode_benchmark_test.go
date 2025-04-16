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
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMultipleLevelRecursiveDep(t *testing.T) {
	// 测试 encode 和 decode 后的对象和源对象结构是一致的，值是相等的
	data := generateLargeMap(2, 10) // about 1M

	encoder := NewEncoder()
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	bytes := encoder.Buffer()

	decoder := NewDecoder(bytes)
	obj, err := decoder.Decode()
	if err != nil {
		panic(err)
	}

	s1 := fmt.Sprintf("%v", obj)
	s2 := fmt.Sprintf("%v", data)
	if s1 != s2 {
		t.Error("deserialize mismatched")
	}
}

func TestMultipleLevelRecursiveDep2(t *testing.T) {
	data := generateLargeMap(3, 5) // about 10MB

	now := time.Now()
	
	encoder := NewEncoder()
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	bytes := encoder.Buffer()
	fmt.Printf("hessian2 serialize %s %dKB\n", time.Since(now), len(bytes)/1024)

	now = time.Now()
	decoder := NewDecoder(bytes)
	obj, err := decoder.Decode()
	if err != nil {
		panic(err)
	}
	rt := time.Since(now)
	fmt.Printf("hessian2 deserialize %s\n", rt)

	if rt > 1*time.Second {
		t.Error("deserialize too slow")
	}
	s1 := fmt.Sprintf("%v", obj)
	s2 := fmt.Sprintf("%v", data)
	if s1 != s2 {
		t.Error("deserialize mismatched")
	}
}
func BenchmarkMultipleLevelRecursiveDep(b *testing.B) {
	// 构造一个多层循环依赖的对象, 无需大对象， 压测情况大批量请求可以看出是否有性能改进
	data := generateLargeMap(2, 5) // about 300KB

	b.ResetTimer()
	for range b.N {
		encoder := NewEncoder()
		err := encoder.Encode(data)
		if err != nil {
			panic(err)
		}
		bytes := encoder.Buffer()

		decoder := NewDecoder(bytes)
		_, err = decoder.Decode()
		if err != nil {
			panic(err)
		}
	}
}

func generateLargeMap(depth int, size int) map[string]interface{} {
	data := map[string]interface{}{}

	if depth != 0 {
		// generate sub map
		for i := 0; i < size; i++ {
			data[fmt.Sprintf("m%d", i)] = generateLargeMap(depth-1, size)
		}

		// generate sub list
		for i := 0; i < size; i++ {
			var sublist []interface{}
			for j := 0; j < size; j++ {
				sublist = append(sublist, generateLargeMap(depth-1, size))
			}
			data[fmt.Sprintf("l%d", i)] = sublist
		}
	}

	// generate string element
	for i := 0; i < size; i++ {
		data[fmt.Sprintf("s%d", i)] = generateRandomString()
	}
	// generate int element
	for i := 0; i < size; i++ {
		data[fmt.Sprintf("i%d", i)] = rand.Int31()
	}
	// generate float element
	for i := 0; i < size; i++ {
		data[fmt.Sprintf("f%d", i)] = rand.Float32()
	}

	return data
}

func generateRandomString() string {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"[rand.Int31n(20):]
}
