package hessian

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestDecodeLargeData(t *testing.T) {
	largeData := generateLargeMap(3, 5) // about 10MB

	now := time.Now()
	largeDataJson, err := json.Marshal(largeData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json serialize %s %dKB\n", time.Since(now), len(largeDataJson)/1024)
	now = time.Now()

	var data map[string]interface{}
	err = json.Unmarshal(largeDataJson, &data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json deserialize %s\n", time.Since(now))
	now = time.Now()

	encoder := NewEncoder()
	err = encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	bytes := encoder.Buffer()
	fmt.Printf("hessian2 serialize %s %dKB\n", time.Since(now), len(bytes)/1024)
	now = time.Now()

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
	if fmt.Sprintf("%v", obj) != fmt.Sprintf("%v", data) {
		t.Error("deserialize mismatched")
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
