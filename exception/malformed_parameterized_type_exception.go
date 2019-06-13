package exception

import (
	hessian "github.com/dubbogo/hessian2"
)
func init(){
	hessian.RegisterPOJO(&MalformedParameterizedTypeException{})
}
type MalformedParameterizedTypeException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []hessian.Exception
	StackTrace           []hessian.StackTraceElement
	Cause                *MalformedParameterizedTypeException
}
func (e MalformedParameterizedTypeException) Error() string {
	return "MalformedParameterizedType"
}

func (MalformedParameterizedTypeException) JavaClassName() string {
	return "java.lang.reflect.MalformedParameterizedTypeException"
}
func NewMalformedParameterizedTypeException(detailMessage string) *MalformedParameterizedTypeException {
	return &MalformedParameterizedTypeException{DetailMessage: detailMessage}
}

