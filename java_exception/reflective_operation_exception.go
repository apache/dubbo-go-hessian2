package java_exception

type ReflectiveOperationException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewReflectiveOperationException(detailMessage string) *ReflectiveOperationException {
	return &ReflectiveOperationException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e ReflectiveOperationException) Error() string {
	return e.DetailMessage
}

func (ReflectiveOperationException) JavaClassName() string {
	return "java.lang.ReflectiveOperationException"
}
