/*
 *
 *  * Copyright 2012-2016 Viant.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 *  * use this file except in compliance with the License. You may obtain a copy of
 *  * the License at
 *  *
 *  * http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 *  * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *  * License for the specific language governing permissions and limitations under
 *  * the License.
 *
 */

// Copyright (c) 2016 ~ 2019, Alex Stocks.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hessian

import (
	jerrors "github.com/juju/errors"
)

const (
	mask = byte(127)
	flag = byte(128)
)

const (
	TAG_READ        = int32(-1)
	ASCII_GAP       = 32
	CHUNK_SIZE      = 4096
	BC_BINARY       = byte('B') // final chunk
	BC_BINARY_CHUNK = byte('A') // non-final chunk

	BC_BINARY_DIRECT  = byte(0x20) // 1-byte length binary
	BINARY_DIRECT_MAX = byte(0x0f)
	BC_BINARY_SHORT   = byte(0x34) // 2-byte length binary
	BINARY_SHORT_MAX  = 0x3ff      // 0-1023 binary

	BC_DATE        = byte(0x4a) // 64-bit millisecond UTC date
	BC_DATE_MINUTE = byte(0x4b) // 32-bit minute UTC date

	BC_DOUBLE = byte('D') // IEEE 64-bit double

	BC_DOUBLE_ZERO  = byte(0x5b)
	BC_DOUBLE_ONE   = byte(0x5c)
	BC_DOUBLE_BYTE  = byte(0x5d)
	BC_DOUBLE_SHORT = byte(0x5e)
	BC_DOUBLE_MILL  = byte(0x5f)

	BC_FALSE = byte('F') // boolean false

	BC_INT = byte('I') // 32-bit int

	INT_DIRECT_MIN = -0x10
	INT_DIRECT_MAX = byte(0x2f)
	BC_INT_ZERO    = byte(0x90)

	INT_BYTE_MIN     = -0x800
	INT_BYTE_MAX     = 0x7ff
	BC_INT_BYTE_ZERO = byte(0xc8)

	BC_END = byte('Z')

	INT_SHORT_MIN     = -0x40000
	INT_SHORT_MAX     = 0x3ffff
	BC_INT_SHORT_ZERO = byte(0xd4)

	BC_LIST_VARIABLE           = byte(0x55)
	BC_LIST_FIXED              = byte('V')
	BC_LIST_VARIABLE_UNTYPED   = byte(0x57)
	BC_LIST_FIXED_UNTYPED      = byte(0x58)
	_listFixedTypedLenTagMin   = byte(0x70)
	_listFixedTypedLenTagMax   = byte(0x77)
	_listFixedUntypedLenTagMin = byte(0x78)
	_listFixedUntypedLenTagMax = byte(0x7f)

	BC_LIST_DIRECT         = byte(0x70)
	BC_LIST_DIRECT_UNTYPED = byte(0x78)
	LIST_DIRECT_MAX        = byte(0x7)

	BC_LONG         = byte('L') // 64-bit signed integer
	LONG_DIRECT_MIN = -0x08
	LONG_DIRECT_MAX = byte(0x0f)
	BC_LONG_ZERO    = byte(0xe0)

	LONG_BYTE_MIN     = -0x800
	LONG_BYTE_MAX     = 0x7ff
	BC_LONG_BYTE_ZERO = byte(0xf8)

	LONG_SHORT_MIN     = -0x40000
	LONG_SHORT_MAX     = 0x3ffff
	BC_LONG_SHORT_ZERO = byte(0x3c)

	BC_LONG_INT = byte(0x59)

	BC_MAP         = byte('M')
	BC_MAP_UNTYPED = byte('H')

	BC_NULL = byte('N') // x4e

	BC_OBJECT     = byte('O')
	BC_OBJECT_DEF = byte('C')

	BC_OBJECT_DIRECT  = byte(0x60)
	OBJECT_DIRECT_MAX = byte(0x0f)

	BC_REF = byte(0x51)

	BC_STRING       = byte('S') // final string
	BC_STRING_CHUNK = byte('R') // non-final string

	BC_STRING_DIRECT  = byte(0x00)
	STRING_DIRECT_MAX = byte(0x1f)
	BC_STRING_SHORT   = byte(0x30)
	STRING_SHORT_MAX  = 0x3ff

	BC_TRUE = byte('T')

	P_PACKET_CHUNK = byte(0x4f)
	P_PACKET       = byte('P')

	P_PACKET_DIRECT   = byte(0x80)
	PACKET_DIRECT_MAX = byte(0x7f)

	P_PACKET_SHORT   = byte(0x70)
	PACKET_SHORT_MAX = 0xfff
	ARRAY_STRING     = "[string"
	ARRAY_INT        = "[int"
	ARRAY_DOUBLE     = "[double"
	ARRAY_FLOAT      = "[float"
	ARRAY_BOOL       = "[boolean"
	ARRAY_LONG       = "[long"

	PATH_KEY      = "path"
	INTERFACE_KEY = "interface"
	VERSION_KEY   = "version"
	TIMEOUT_KEY   = "timeout"

	STRING_NIL   = "null"
	STRING_TRUE  = "true"
	STRING_FALSE = "false"
	STRING_ZERO  = "0.0"
	STRING_ONE   = "1.0"
)

const (
	Response_OK                byte = 20
	Response_CLIENT_TIMEOUT    byte = 30
	Response_SERVER_TIMEOUT    byte = 31
	Response_BAD_REQUEST       byte = 40
	Response_BAD_RESPONSE      byte = 50
	Response_SERVICE_NOT_FOUND byte = 60
	Response_SERVICE_ERROR     byte = 70
	Response_SERVER_ERROR      byte = 80
	Response_CLIENT_ERROR      byte = 90

	RESPONSE_WITH_EXCEPTION                  int32 = 0
	RESPONSE_VALUE                           int32 = 1
	RESPONSE_NULL_VALUE                      int32 = 2
	RESPONSE_WITH_EXCEPTION_WITH_ATTACHMENTS int32 = 3
	RESPONSE_VALUE_WITH_ATTACHMENTS          int32 = 4
	RESPONSE_NULL_VALUE_WITH_ATTACHMENTS     int32 = 5
)

/**
 * 协议头是16字节的定长数据
 * 2字节magic字符串0xdabb,0-7高位，8-15低位
 * 1字节的消息标志位。16-20序列id,21 event,22 two way,23请求或响应标识
 * 1字节状态。当消息类型为响应时，设置响应状态。24-31位。
 * 8字节，消息ID,long类型，32-95位。
 * 4字节，消息长度，96-127位
 **/
const (
	// header length.
	HEADER_LENGTH = 16

	// magic header
	MAGIC      = uint16(0xdabb)
	MAGIC_HIGH = byte(0xda)
	MAGIC_LOW  = byte(0xbb)

	// message flag.
	FLAG_REQUEST = byte(0x80)
	FLAG_TWOWAY  = byte(0x40)
	FLAG_EVENT   = byte(0x20) // for heartbeat
	SERIAL_MASK  = 0x1f

	DUBBO_VERSION                          = "2.5.4"
	DUBBO_VERSION_KEY                      = "dubbo"
	DEFAULT_DUBBO_PROTOCOL_VERSION         = "2.0.2" // Dubbo RPC protocol version, for compatibility, it must not be between 2.0.10 ~ 2.6.2
	LOWEST_VERSION_FOR_RESPONSE_ATTACHMENT = 2000200
	DEFAULT_LEN                            = 8388608 // 8 * 1024 * 1024 default body max length
)

var (
	DubboRequestHeaderBytes      = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY}
	DubboResponseHeaderBytes     = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, 0x00, Response_OK}
	DubboRequestHeartbeatHeader  = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY | FLAG_EVENT}
	DubboResponseHeartbeatHeader = [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_EVENT}
)

var (
	ErrHeaderNotEnough = jerrors.New("header buffer too short")
	ErrBodyNotEnough   = jerrors.New("body buffer too short")
	ErrJavaException   = jerrors.New("got java exception")
	ErrIllegalPackage  = jerrors.New("illegal package!")
)
