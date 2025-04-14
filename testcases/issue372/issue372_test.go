package issue372

import (
	"bufio"
	"bytes"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MediaStream struct {
	Payload string
}

func (MediaStream) JavaClassName() string {
	return "com.test.MediaStream"
}

func TestDecodeFromTcpStream(t *testing.T) {
	payload := make([]byte, 1024)
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for i, _ := range payload {
		payload[i] = alphabet[i%26]
	}
	info := &MediaStream{
		Payload: string(payload),
	}

	hessian.RegisterPOJO(info)
	codecW := hessian.NewHessianCodec(nil)
	service := hessian.Service{
		Path:      "test",
		Interface: "ITest",
		Version:   "v1.0",
		Method:    "test",
		Timeout:   time.Second * 10,
	}
	header := hessian.DubboHeader{
		SerialID:       2,
		Type:           hessian.PackageRequest,
		ID:             1,
		ResponseStatus: hessian.Zero,
	}
	resp, err := codecW.Write(service, header, []interface{}{info})

	// set reader buffer = 1024 to split resp into two parts
	codec := hessian.NewHessianCodecCustom(0, bufio.NewReaderSize(bytes.NewReader(resp), 1024), 0, true)
	h := &hessian.DubboHeader{}
	assert.NoError(t, codec.ReadHeader(h))
	assert.Equal(t, h.SerialID, header.SerialID)
	assert.Equal(t, h.Type, header.Type)
	assert.Equal(t, h.ID, header.ID)
	assert.Equal(t, h.ResponseStatus, header.ResponseStatus)

	reqBody := make([]interface{}, 7)

	err = codec.ReadBody(reqBody)
	assert.NoError(t, err)
	assert.Equal(t, reqBody[1], service.Path)
	assert.Equal(t, reqBody[2], service.Version)
	assert.Equal(t, reqBody[3], service.Method)

	if list, ok := reqBody[5].([]interface{}); ok {
		if infoPtr, ok2 := list[0].(*MediaStream); ok2 {
			assert.Equal(t, len(infoPtr.Payload), 1024)
		}
	}

	codec = hessian.NewHessianCodecCustom(0, bufio.NewReaderSize(bytes.NewReader(resp), 1024), 0, false)
	assert.Error(t, codec.ReadHeader(h))
}
