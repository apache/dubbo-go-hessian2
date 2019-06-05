package hessian

func init() {
	RegisterPOJO(&Throwable{})
	RegisterPOJO(&Exception{})
	RegisterPOJO(&StackTraceElement{})
}

////////////////////////////
// Throwable interface
////////////////////////////

type ThrowableIntf interface {
	Error() string
	JavaClassName() string
}

////////////////////////////
// Throwable
////////////////////////////

type Throwable struct {
	SerialVersionUID int64
	DetailMessage    string
	//todo:backtrace
	SuppressedExceptions []Throwable
	StackTrace           []StackTraceElement
	Cause                *Throwable
}

func NewThrowable(detailMessage string) *Throwable {
	return &Throwable{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e Throwable) Error() string {
	return e.DetailMessage
}

func (Throwable) JavaClassName() string {
	return "java.lang.Throwable"
}

////////////////////////////
// Exception
////////////////////////////

type Exception struct {
	SerialVersionUID int64
	DetailMessage    string
	//todo:backtrace
	SuppressedExceptions []Exception
	StackTrace           []StackTraceElement
	Cause                *Exception
}

func NewException(detailMessage string) *Exception {
	return &Exception{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e Exception) Error() string {
	return e.DetailMessage
}

func (Exception) JavaClassName() string {
	return "java.lang.Exception"
}

type StackTraceElement struct {
	DeclaringClass string
	MethodName     string
	FileName       string
	LineNumber     int
}

func (StackTraceElement) JavaClassName() string {
	return "java.lang.StackTraceElement"
}
