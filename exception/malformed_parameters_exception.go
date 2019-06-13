package exception


import (
	hessian "github.com/dubbogo/hessian2"
)
func init(){
	hessian.RegisterPOJO(&MalformedParametersException{})
}
type MalformedParametersException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []hessian.Exception
	StackTrace           []hessian.StackTraceElement
	Cause                *MalformedParametersException
}
func (e MalformedParametersException) Error() string {
	return e.DetailMessage
}

func (MalformedParametersException) JavaClassName() string {
	return "java.lang.reflect.MalformedParametersException"
}
func NewMalformedParametersException(detailMessage string) *MalformedParametersException {
	return &MalformedParametersException{DetailMessage: detailMessage}
}

