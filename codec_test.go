/******************************************************
# DESC    : codec.go unittest
# AUTHOR  : Alex Stocks
# EMAIL   : alexstocks@foxmail.com
# MOD     : 2016-10-23 19:51
# FILE    : codec_test.go
******************************************************/

package hessian

import (
	"testing"
)

// go test -v  codec_test.go encode.go   const.go  pojo.go codec.go
// go test -v -run TestPackUint16

func TestPackUint16(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackUint16(PackUint16(uint16(0xfedc), arr)))
	var v = uint16(0xfedc)
	if r := UnpackUint16(PackUint16(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt16(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt16(PackInt16(int16(0x1234), arr)))
	var v = int16(0x1234)
	if r := UnpackInt16(PackInt16(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt32(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt32(PackInt32(int32(0x12344678), arr)))
	var v = int32(0x12344678)
	if r := UnpackInt32(PackInt32(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}

func TestPackInt64(t *testing.T) {
	// var arr []byte
	// t.Logf("0X%x\n", UnpackInt64(PackInt64(int64(0x1234567890abcdef), arr)))
	var v = int64(0x1234567890abcdef)
	if r := UnpackInt64(PackInt64(v)); r != v {
		t.Fatalf("v:0X%d, pack-unpack value:0X%x\n", v, r)
	}
}
