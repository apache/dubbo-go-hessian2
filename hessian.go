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
	"bufio"
	"encoding/binary"
	"time"
)

import (
	jerrors "github.com/juju/errors"
)

const (
	Error          PackgeType = 0x01
	Request                   = 0x02
	Response                  = 0x04
	Heartbeat                 = 0x08
	Request_TwoWay            = 0x10
)

type PackgeType int

type DubboHeader struct {
	SerialID       byte
	Type           PackgeType
	ID             int64
	BodyLen        int
	ResponseStatus byte
}

type Service struct {
	Path      string
	Interface string
	Version   string
	Target    string // Service Name
	Method    string
	Timeout   time.Duration // request timeout
}

type HessianCodec struct {
	pkgType PackgeType
	reader  *bufio.Reader
	bodyLen int
}

func NewHessianCodec(reader *bufio.Reader) *HessianCodec {
	return &HessianCodec{
		reader: reader,
	}
}

func (h *HessianCodec) Write(service Service, header DubboHeader, body interface{}) ([]byte, error) {
	switch header.Type {
	case Heartbeat:
		if header.ResponseStatus == Zero {
			return PackRequest(service, header, body)
		}
		return PackResponse(header, map[string]string{}, body)
	case Request:
		return PackRequest(service, header, body)

	case Response:
		return PackResponse(header, map[string]string{}, body)

	default:
		return nil, jerrors.Errorf("Unrecognised message type: %v", header.Type)
	}

	return nil, nil
}

func (h *HessianCodec) ReadHeader(header *DubboHeader) error {

	var err error

	buf, err := h.reader.Peek(HEADER_LENGTH)
	if err != nil { // this is impossible
		return jerrors.Trace(err)
	}
	_, err = h.reader.Discard(HEADER_LENGTH)
	if err != nil { // this is impossible
		return jerrors.Trace(err)
	}

	//// read header

	if buf[0] != byte(MAGIC_HIGH) && buf[1] != byte(MAGIC_LOW) {
		return ErrIllegalPackage
	}

	// Header{serialization id(5 bit), event, two way, req/response}
	if header.SerialID = buf[2] & SERIAL_MASK; header.SerialID == Zero {
		return jerrors.Errorf("serialization ID:%v", header.SerialID)
	}

	flag := buf[2] & FLAG_EVENT
	if flag != Zero {
		header.Type |= Heartbeat
	}
	flag = buf[2] & FLAG_REQUEST
	if flag != Zero {
		header.Type |= Request
		flag = buf[2] & FLAG_TWOWAY
		if flag != Zero {
			header.Type |= Request_TwoWay
		}
	} else {
		header.Type |= Response

		// Header{status}
		if buf[3] != Response_OK {
			err = ErrJavaException
			header.Type |= Error
			bufSize := h.reader.Buffered()
			if bufSize > 2 { // responseType + objectType + error content,so it's size > 2
				expBuf, expErr := h.reader.Peek(bufSize)
				if expErr == nil {
					err = jerrors.Errorf("java exception:%s", string(expBuf[2:bufSize-1]))
				}
			}
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

	return jerrors.Trace(err)

}

func (h *HessianCodec) ReadBody(rspObj interface{}) error {

	buf, err := h.reader.Peek(h.bodyLen)
	if err == bufio.ErrBufferFull {
		return ErrBodyNotEnough
	}
	if err != nil {
		return jerrors.Trace(err)
	}
	_, err = h.reader.Discard(h.bodyLen)
	if err != nil { // this is impossible
		return jerrors.Trace(err)
	}

	switch h.pkgType & 0x0f {
	case Request | Heartbeat, Response | Heartbeat:
		return nil
	case Request:
		if rspObj != nil {
			if err = UnpackRequestBody(buf, rspObj); err != nil {
				return jerrors.Trace(err)
			}
		}

		return nil

	case Response:
		if rspObj != nil {
			if err = UnpackResponseBody(buf, rspObj); err != nil {
				return jerrors.Trace(err)
			}
		}
	}

	return nil
}
