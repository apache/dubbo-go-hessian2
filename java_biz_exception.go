package hessian

import (
	"fmt"
	"github.com/apache/dubbo-go-hessian2/java_exception"
	"sync"
)

var mutex sync.Mutex

func checkAndGetException(cls classInfo) (structInfo, bool){

	if len(cls.fieldNameList) < 4 {
		return structInfo{}, false
	}
	var (
		throwable structInfo
		ok bool
	)
	var count =0
	for _, item := range cls.fieldNameList {
		if item == "detailMessage" || item == "suppressedExceptions" || item == "stackTrace" || item == "cause" {
			count ++
		}
	}
	// 如果满足异常条件
	if count == 4 {
		mutex.Lock()
		defer mutex.Unlock()
		if throwable, ok = getStructInfo(cls.javaName); ok {
			return throwable, true
		}
		RegisterPOJO(newBizException(cls.javaName))
		if throwable, ok = getStructInfo(cls.javaName); ok {
			return throwable, true
		}
	}
	return throwable, count == 4
}

type BizException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []java_exception.Throwabler
	StackTrace           []java_exception.StackTraceElement
	Cause                java_exception.Throwabler
	name				 string
}

// NewThrowable is the constructor
func newBizException(name string) *BizException {
	return &BizException{name: name, StackTrace: []java_exception.StackTraceElement{}}
}

// Error output error message
func (e BizException) Error() string {
	return fmt.Sprintf("throw %v : %v", e.name, e.DetailMessage)
}

//JavaClassName  java fully qualified path
func (e BizException) JavaClassName() string {
	return e.name
}