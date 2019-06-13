package exception

import (
	hessian "github.com/dubbogo/hessian2"
)
func init(){
	hessian.RegisterPOJO(&TypeNotPresentException{})
}
type TypeNotPresentException struct {
	TypeName         string
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []hessian.Exception
	StackTrace           []hessian.StackTraceElement
	Cause                *hessian.Throwable

}
func (e TypeNotPresentException) Error() string {
	return e.DetailMessage
}

func (TypeNotPresentException) JavaClassName() string {
	return "java.lang.TypeNotPresentException"
}
func NewTypeNotPresentException(typeName string,detailMessage string) *TypeNotPresentException {
	return &TypeNotPresentException{TypeName: typeName,DetailMessage:detailMessage}
}
