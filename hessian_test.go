package hessian

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

type Case struct {
	A string
	B int
}

func (c *Case) JavaClassName() string {
	return "com.test.case"
}

func doTest(t *testing.T, packageType PackgeType, responseStatus byte, body interface{}) {
	RegisterPOJO(&Case{})

	codecW := NewHessianCodec(nil)
	resp, err := codecW.Write(Service{
		Path:      "/test",
		Interface: "ITest",
		Version:   "v1.0",
		Target:    "test",
		Method:    "test",
		Timeout:   time.Second * 10,
	}, DubboHeader{
		SerialID:       2,
		Type:           packageType,
		ID:             1,
		ResponseStatus: responseStatus,
	}, body)
	assert.Nil(t, err)

	codecR := NewHessianCodec(bufio.NewReader(bytes.NewReader(resp)))

	h := &DubboHeader{}
	err = codecR.ReadHeader(h)
	if responseStatus == Response_OK || responseStatus == Zero {
		assert.Nil(t, err)
	} else {
		t.Log(err)
		assert.NotNil(t, err)
		return
	}
	assert.Equal(t, byte(2), h.SerialID)
	assert.Equal(t, packageType, h.Type&(Request|Response|Heartbeat))
	assert.Equal(t, int64(1), h.ID)
	assert.Equal(t, responseStatus, h.ResponseStatus)

	var c interface{}
	n := reflect.TypeOf(body).String()
	if n == "*hessian.Case" {
		c = &Case{}
	} else if n == "string" {
		tmp := ""
		c = &tmp
		c = &c
	} else if n == "int64" {
		tmp := 1
		c = &tmp
		c = &c
	} else if n == "bool" {
		tmp := false
		c = &tmp
		c = &c
	} else {
		c = make([]interface{}, 7)
	}
	err = codecR.ReadBody(c)
	assert.Nil(t, err)
	t.Log(c)
	//t.Log(reflect.ValueOf(body).Type())
	//t.Log(reflect.ValueOf(c).Type())
	if packageType == Request {
		assert.True(t, len(body.([]interface{})) == len(c.([]interface{})[5].([]interface{})))
	} else if packageType == Response {
		assert.True(t, reflect.DeepEqual(body, c))
	}
}

func TestResponse(t *testing.T) {
	doTest(t, Response, Response_OK, &Case{A: "a", B: 1})
	doTest(t, Response, Response_OK, "ok!!!!!")
	doTest(t, Response, Response_OK, int64(3))
	doTest(t, Response, Response_OK, true)
	doTest(t, Response, Response_SERVER_ERROR, "error!!!!!")
}

func TestRequest(t *testing.T) {
	doTest(t, Request, byte(0), []interface{}{"a"})
	doTest(t, Request, byte(0), []interface{}{"a", 3})
	doTest(t, Request, byte(0), []interface{}{"a", true})
	doTest(t, Request, byte(0), []interface{}{"a", 3, true})
	doTest(t, Request, byte(0), []interface{}{3.2, true})
	doTest(t, Request, byte(0), []interface{}{"a", 3, true, &Case{A: "a", B: 3}})
	doTest(t, Request, byte(0), []interface{}{"a", 3, true, []*Case{{A: "a", B: 3}}})
}
