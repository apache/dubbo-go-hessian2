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

// go test -v -run TestPackUint16
func TestPackUint16(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackUint16(PackUint16(uint16(0xfedc), arr)))
	v := uint16(0xfedc)
	if r := UnpackUint16(PackUint16(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt16(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt16(PackInt16(int16(0x1234), arr)))
	v := int16(0x1234)
	if r := UnpackInt16(PackInt16(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt32(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt32(PackInt32(int32(0x12344678), arr)))
	v := int32(0x12344678)
	if r := UnpackInt32(PackInt32(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt64(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt64(PackInt64(int64(0x1234567890abcdef), arr)))
	v := int64(0x1234567890abcdef)
	if r := UnpackInt64(PackInt64(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}
