package hessian

func init() {
	RegisterPOJO(&Exception{})
	RegisterPOJO(&StackTraceElement{})
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
	LineNumber     string
}

func (StackTraceElement) JavaClassName() string {
	return "java.lang.StackTraceElement"
}
