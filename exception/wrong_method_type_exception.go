package exception

import (
hessian "github.com/dubbogo/hessian2"
)
func init(){
	hessian.RegisterPOJO(&WrongMethodTypeException{})
}
type WrongMethodTypeException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []hessian.Exception
	StackTrace           []hessian.StackTraceElement
	Cause                *WrongMethodTypeException
}
func (e WrongMethodTypeException) Error() string {
	return e.DetailMessage
}

func (WrongMethodTypeException) JavaClassName() string {
	return "java.lang.invoke.WrongMethodTypeException"
}
func NewWrongMethodTypeException(detailMessage string) *WrongMethodTypeException {
	return &WrongMethodTypeException{DetailMessage: detailMessage}
}
