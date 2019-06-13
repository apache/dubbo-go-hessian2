
package exception

import (
hessian "github.com/dubbogo/hessian2"
)
func init(){
	hessian.RegisterPOJO(&UndeclaredThrowableException{})
}
type UndeclaredThrowableException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []hessian.Exception
	StackTrace           []hessian.StackTraceElement
	Cause                *hessian.Throwable
	UndeclaredThrowable        hessian.Throwable
}
func (e UndeclaredThrowableException) Error() string {
	return e.DetailMessage
}

func (UndeclaredThrowableException) JavaClassName() string {
	return "java.lang.reflect.UndeclaredThrowableException"
}
func NewUndeclaredThrowableException(detailMessage string) *UndeclaredThrowableException {
	return &UndeclaredThrowableException{DetailMessage: detailMessage,UndeclaredThrowable:hessian.Throwable{}}
}
