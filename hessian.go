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

// This file contains the implementation of
// [dubbo2 rpc protocol](https://dubbo.apache.org/zh/docs/concepts/rpc-protocol/#dubbo2), which is deprecated.
// Use the [implementation](https://github.com/apache/dubbo-go/tree/master/protocol/dubbo/hessian2) in dubbo-go project instead.

package hessian

import (
	"bufio"
	"encoding/binary"
	"time"
)

import (
	perrors "github.com/pkg/errors"
)

// enum part
const (
	PackageError              = PackageType(0x01)
	PackageRequest            = PackageType(0x02)
	PackageResponse           = PackageType(0x04)
	PackageHeartbeat          = PackageType(0x08)
	PackageRequest_TwoWay     = PackageType(0x10)
	PackageResponse_Exception = PackageType(0x20)
	PackageType_BitSize       = 0x2f
)

// PackageType defines dubbo package type.
type PackageType int

// DubboHeader dubbo header
type DubboHeader struct {
	SerialID       byte
	Type           PackageType
	ID             int64
	BodyLen        int
	ResponseStatus byte
}

// Service defines service instance
type Service struct {
	Path      string
	Interface string
	Group     string
	Version   string
	Method    string
	Timeout   time.Duration // request timeout
}

// HessianCodec defines hessian codec
type HessianCodec struct {
	pkgType  PackageType
	reader   *bufio.Reader
	bodyLen  int
	isStream bool
}

// NewHessianCodec generate a new hessian codec instance
func NewHessianCodec(reader *bufio.Reader) *HessianCodec {
	return &HessianCodec{
		reader: reader,
	}
}

// NewHessianCodecCustom generate a new hessian codec instance.
func NewHessianCodecCustom(pkgType PackageType, reader *bufio.Reader, bodyLen int, isStream bool) *HessianCodec {
	return &HessianCodec{
		pkgType:  pkgType,
		reader:   reader,
		bodyLen:  bodyLen,
		isStream: isStream,
	}
}

func (h *HessianCodec) Write(service Service, header DubboHeader, body interface{}) ([]byte, error) {
	switch header.Type {
	case PackageHeartbeat:
		if header.ResponseStatus == Zero {
			return packRequest(service, header, body)
		}
		return packResponse(header, body)

	case PackageRequest, PackageRequest_TwoWay:
		return packRequest(service, header, body)

	case PackageResponse:
		return packResponse(header, body)

	default:
		return nil, perrors.Errorf("Unrecognized message type: %v", header.Type)
	}
}

// ReadHeader uses hessian codec to read dubbo header
func (h *HessianCodec) ReadHeader(header *DubboHeader) error {
	var err error

	var buf []byte

	if !h.isStream {
		if h.reader.Size() < HEADER_LENGTH {
			return ErrHeaderNotEnough
		}
		buf, err = h.reader.Peek(HEADER_LENGTH)
		if err != nil { // this is impossible
			return perrors.WithStack(err)
		}
		_, err = h.reader.Discard(HEADER_LENGTH)
		if err != nil { // this is impossible
			return perrors.WithStack(err)
		}
	} else {
		buf = make([]byte, HEADER_LENGTH)
		n, _ := h.reader.Read(buf)
		if n < HEADER_LENGTH {
			// just try once
			_, err = h.reader.Peek(HEADER_LENGTH - n)
			if err != nil {
				return perrors.WithStack(err)
			}
			_, err = h.reader.Read(buf[n:])
			if err != nil { // this is impossible
				return perrors.WithStack(err)
			}
		}
	}

	//// read header

	if buf[0] != MAGIC_HIGH && buf[1] != MAGIC_LOW {
		return ErrIllegalPackage
	}

	// Header{serialization id(5 bit), event, two way, req/response}
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == Zero {
		return perrors.Errorf("serialization ID:%v", header.SerialID)
	}

	headerFlag := buf[2] & FLAG_EVENT
	if headerFlag != Zero {
		header.Type |= PackageHeartbeat
	}
	headerFlag = buf[2] & FLAG_REQUEST
	if headerFlag != Zero {
		header.Type |= PackageRequest
		headerFlag = buf[2] & FLAG_TWOWAY
		if headerFlag != Zero {
			header.Type |= PackageRequest_TwoWay
		}
	} else {
		header.Type |= PackageResponse
		header.ResponseStatus = buf[3]
		if header.ResponseStatus != Response_OK {
			header.Type |= PackageResponse_Exception
		}
	}

	// Header{req id}
	header.ID = int64(binary.BigEndian.Uint64(buf[4:]))

	// Header{body len}
	header.BodyLen = int(binary.BigEndian.Uint32(buf[12:]))
	if header.BodyLen < 0 {
		return ErrIllegalPackage
	}

	h.pkgType = header.Type
	h.bodyLen = header.BodyLen

	if h.reader.Buffered() < h.bodyLen && !h.isStream {
		return ErrBodyNotEnough
	}

	return perrors.WithStack(err)
}

// ReadBody uses hessian codec to read response body
func (h *HessianCodec) ReadBody(rspObj interface{}) error {
	var err error
	var buf []byte
	if !h.isStream {
		if h.reader.Buffered() < h.bodyLen {
			return ErrBodyNotEnough
		}
		buf, err = h.reader.Peek(h.bodyLen)
		if err != nil {
			return perrors.WithStack(err)
		}
		_, err = h.reader.Discard(h.bodyLen)
		if err != nil { // this is impossible
			return perrors.WithStack(err)
		}
	} else {
		buf = make([]byte, h.bodyLen)
		readLen, n := 0, 0
		for readLen < h.bodyLen {
			n, err = h.reader.Read(buf[readLen:])
			if err != nil {
				return perrors.WithStack(err)
			}
			readLen += n
		}
	}

	switch h.pkgType & PackageType_BitSize {
	case PackageResponse | PackageHeartbeat | PackageResponse_Exception, PackageResponse | PackageResponse_Exception:
		decoder := NewDecoder(buf[:])
		exception, decErr := decoder.Decode()
		if decErr != nil {
			return perrors.WithStack(decErr)
		}
		rsp, ok := rspObj.(*Response)
		if !ok {
			return perrors.Errorf("java exception:%s", exception.(string))
		}
		rsp.Exception = perrors.Errorf("java exception:%s", exception.(string))
		return nil
	case PackageRequest | PackageHeartbeat, PackageResponse | PackageHeartbeat:
	case PackageRequest:
		if rspObj != nil {
			if err = unpackRequestBody(NewStrictDecoder(buf[:]), rspObj); err != nil {
				return perrors.WithStack(err)
			}
		}
	case PackageResponse:
		if rspObj != nil {
			if err = unpackResponseBody(NewDecoder(buf[:]), rspObj); err != nil {
				return perrors.WithStack(err)
			}
		}
	}

	return nil
}

// ReadAttachments ignore body, but only read attachments
func (h *HessianCodec) ReadAttachments() (map[string]string, error) {
	var err error
	var buf []byte
	if !h.isStream {
		if h.reader.Buffered() < h.bodyLen {
			return nil, ErrBodyNotEnough
		}
		buf, err = h.reader.Peek(h.bodyLen)
		if err != nil {
			return nil, perrors.WithStack(err)
		}
		_, err = h.reader.Discard(h.bodyLen)
		if err != nil { // this is impossible
			return nil, perrors.WithStack(err)
		}
	} else {
		buf = make([]byte, h.bodyLen)
		readLen, n := 0, 0
		for readLen < h.bodyLen {
			n, err = h.reader.Read(buf[readLen:])
			if err != nil {
				return nil, perrors.WithStack(err)
			}
			readLen += n
		}
	}

	switch h.pkgType & PackageType_BitSize {
	case PackageRequest:
		rspObj := make([]interface{}, 7)
		if err = unpackRequestBody(NewDecoderWithSkip(buf[:]), rspObj); err != nil {
			return nil, perrors.WithStack(err)
		}
		return rspObj[6].(map[string]string), nil
	case PackageResponse:
		rspObj := &Response{}
		if err = unpackResponseBody(NewDecoderWithSkip(buf[:]), rspObj); err != nil {
			return nil, perrors.WithStack(err)
		}
		return rspObj.Attachments, nil
	}

	return nil, nil
}
