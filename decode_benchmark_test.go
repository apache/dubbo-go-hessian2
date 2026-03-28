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
	"math"
	"math/rand"
	"testing"
	"time"
)

// TestMultipleLevelRecursiveDep verifies that encode followed by decode
// produces a value-equivalent nested map without relying on map iteration order.
func TestMultipleLevelRecursiveDep(t *testing.T) {
	data := generateLargeMap(2, 8) // ~500KB nested map

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

	if err := compareDecodedValue("root", data, obj); err != nil {
		t.Fatal(err)
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

	if obj == nil {
		b.Fatal("deserialize result is nil")
	}

	if _, ok := obj.(map[interface{}]interface{}); !ok {
		b.Fatalf("deserialize result type mismatch, expected map, got %T", obj)
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

// BenchmarkMultipleLevelRecursiveDep benchmarks a full encode+decode cycle
// on a medium-sized (~300KB) nested map.
func BenchmarkMultipleLevelRecursiveDep(b *testing.B) {
	data := generateLargeMap(2, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder := NewEncoder()
		err := encoder.Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		bytes := encoder.Buffer()

		decoder := NewDecoder(bytes)
		_, err = decoder.Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// generateLargeMap builds a nested map with sub-maps, sub-lists, strings,
// ints, and floats. depth controls nesting levels; size controls fan-out.
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

func compareDecodedValue(path string, original interface{}, decoded interface{}) error {
	switch originalValue := original.(type) {
	case map[string]interface{}:
		return compareDecodedMap(path, originalValue, decoded)
	case []interface{}:
		return compareDecodedSlice(path, originalValue, decoded)
	case float32:
		decodedFloat, err := decodedFloat64(path, decoded)
		if err != nil {
			return err
		}
		if math.Abs(float64(originalValue)-decodedFloat) > 1e-6 {
			return fmt.Errorf("%s: float mismatch, expected %f, got %f", path, originalValue, decodedFloat)
		}
	case string:
		decodedString, ok := decoded.(string)
		if !ok {
			return fmt.Errorf("%s: type mismatch, expected string, got %T", path, decoded)
		}
		if originalValue != decodedString {
			return fmt.Errorf("%s: string mismatch, expected %q, got %q", path, originalValue, decodedString)
		}
	case int32:
		decodedInt, err := decodedInt64(path, decoded)
		if err != nil {
			return err
		}
		if int64(originalValue) != decodedInt {
			return fmt.Errorf("%s: int mismatch, expected %d, got %d", path, originalValue, decodedInt)
		}
	default:
		if original != decoded {
			return fmt.Errorf("%s: value mismatch, expected %v, got %v", path, original, decoded)
		}
	}

	return nil
}

func compareDecodedMap(path string, original map[string]interface{}, decoded interface{}) error {
	decodedMap, ok := decoded.(map[interface{}]interface{})
	if !ok {
		return fmt.Errorf("%s: type mismatch, expected map, got %T", path, decoded)
	}
	if len(original) != len(decodedMap) {
		return fmt.Errorf("%s: map size mismatch, expected %d, got %d", path, len(original), len(decodedMap))
	}

	for key, originalValue := range original {
		decodedValue, exists := decodedMap[key]
		if !exists {
			return fmt.Errorf("%s.%s: key missing in decoded map", path, key)
		}
		if err := compareDecodedValue(path+"."+key, originalValue, decodedValue); err != nil {
			return err
		}
	}

	return nil
}

func compareDecodedSlice(path string, original []interface{}, decoded interface{}) error {
	decodedSlice, ok := decoded.([]interface{})
	if !ok {
		return fmt.Errorf("%s: type mismatch, expected slice, got %T", path, decoded)
	}
	if len(original) != len(decodedSlice) {
		return fmt.Errorf("%s: slice length mismatch, expected %d, got %d", path, len(original), len(decodedSlice))
	}

	for i, originalValue := range original {
		if err := compareDecodedValue(fmt.Sprintf("%s[%d]", path, i), originalValue, decodedSlice[i]); err != nil {
			return err
		}
	}

	return nil
}

func decodedFloat64(path string, value interface{}) (float64, error) {
	switch decodedValue := value.(type) {
	case float32:
		return float64(decodedValue), nil
	case float64:
		return decodedValue, nil
	default:
		return 0, fmt.Errorf("%s: type mismatch, expected float, got %T", path, value)
	}
}

func decodedInt64(path string, value interface{}) (int64, error) {
	switch decodedValue := value.(type) {
	case int:
		return int64(decodedValue), nil
	case int8:
		return int64(decodedValue), nil
	case int16:
		return int64(decodedValue), nil
	case int32:
		return int64(decodedValue), nil
	case int64:
		return decodedValue, nil
	default:
		return 0, fmt.Errorf("%s: type mismatch, expected int, got %T", path, value)
	}
}
