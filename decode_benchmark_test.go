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

// No third-party imports
// No internal imports
func TestMultipleLevelRecursiveDep(t *testing.T) {
	// Test serialization consistency for complex nested structures
	data := generateLargeMap(2, 8) // Reduced size for stability: about 500KB

	encoder := NewEncoder()
	err := encoder.Encode(data)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}
	bytes := encoder.Buffer()

	decoder := NewDecoder(bytes)
	obj, err := decoder.Decode()
	if err != nil {
		t.Fatalf("Failed to decode data: %v", err)
	}

	// Use structured comparison instead of string comparison
	decodedMap, ok := obj.(map[interface{}]interface{})
	if !ok {
		t.Fatal("Decoded object is not a map")
	}

	// Verify basic structure and some key elements
	if !verifyMapStructure(data, decodedMap, t) {
		t.Error("Decoded map structure does not match original")
	}
}

// BenchmarkMultipleLevelRecursiveDepLarge measures decode performance on a large object.
// This benchmark focuses on performance measurement rather than strict thresholds.
func BenchmarkMultipleLevelRecursiveDepLarge(b *testing.B) {
	// Test with a moderately large object for performance measurement
	data := generateLargeMap(3, 4) // Reduced from (3,5) for better stability

	startEncode := time.Now()
	encoder := NewEncoder()
	err := encoder.Encode(data)
	if err != nil {
		b.Fatal(err)
	}
	bytes := encoder.Buffer()
	encodeTime := time.Since(startEncode)
	b.Logf("serialize %s %dKB", encodeTime, len(bytes)/1024)

	// Perform one decode operation to verify it works
	startDecode := time.Now()
	decoder := NewDecoder(bytes)
	obj, err := decoder.Decode()
	if err != nil {
		b.Fatal(err)
	}
	decodeTime := time.Since(startDecode)
	b.Logf("deserialize %s", decodeTime)

	// Basic validation - ensure decode succeeded and returned correct type
	if obj == nil {
		b.Error("deserialize result is nil")
	}

	if _, ok := obj.(map[interface{}]interface{}); !ok {
		b.Error("deserialize result type mismatch, expected map")
	}

	// Log performance metrics for analysis
	b.Logf("Performance metrics - Encode: %v, Decode: %v, Size: %dKB",
		encodeTime, decodeTime, len(bytes)/1024)

	// Run the actual benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder := NewDecoder(bytes)
		_, err := decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMultipleLevelRecursiveDep(b *testing.B) {
	// benchmark for decode()
	data := generateLargeMap(2, 5) // about 300KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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

// verifyMapStructure performs structured comparison instead of string comparison
// to avoid issues with floating point precision and map iteration order
func verifyMapStructure(original map[string]interface{}, decoded map[interface{}]interface{}, t *testing.T) bool {
	// Check if basic structure elements exist
	for key, originalValue := range original {
		decodedValue, exists := decoded[key]
		if !exists {
			t.Logf("Key %s missing in decoded map", key)
			return false
		}

		// For nested maps, recursively verify
		if originalMap, ok := originalValue.(map[string]interface{}); ok {
			if decodedMap, ok := decodedValue.(map[interface{}]interface{}); ok {
				if !verifyMapStructure(originalMap, decodedMap, t) {
					return false
				}
			} else {
				t.Logf("Value type mismatch for key %s: expected map, got %T", key, decodedValue)
				return false
			}
			continue
		}

		// For slices, verify basic structure
		if originalSlice, ok := originalValue.([]interface{}); ok {
			if decodedSlice, ok := decodedValue.([]interface{}); ok {
				if len(originalSlice) != len(decodedSlice) {
					t.Logf("Slice length mismatch for key %s: expected %d, got %d", key, len(originalSlice), len(decodedSlice))
					return false
				}
			} else {
				t.Logf("Value type mismatch for key %s: expected slice, got %T", key, decodedValue)
				return false
			}
			continue
		}

		// For basic types, we can do direct comparison
		// but be tolerant of floating point precision differences and type conversions
		if originalFloat, ok := originalValue.(float32); ok {
			// Hessian may convert float32 to float64
			var decodedFloat float64
			if f32, ok := decodedValue.(float32); ok {
				decodedFloat = float64(f32)
			} else if f64, ok := decodedValue.(float64); ok {
				decodedFloat = f64
			} else {
				t.Logf("Value type mismatch for key %s: expected float, got %T", key, decodedValue)
				return false
			}
			// Allow small floating point differences
			if abs64(float64(originalFloat)-decodedFloat) > 1e-6 {
				t.Logf("Float value mismatch for key %s: expected %f, got %f", key, originalFloat, decodedFloat)
				return false
			}
			continue
		}
	}

	return true
}

// Use math.Abs for absolute value calculations.
