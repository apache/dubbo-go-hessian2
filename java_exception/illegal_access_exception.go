package java_exception

type IllegalAccessException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewIllegalAccessException(detailMessage string) *IllegalAccessException {
	return &IllegalAccessException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e IllegalAccessException) Error() string {
	return e.DetailMessage
}

func (IllegalAccessException) JavaClassName() string {
	return "java.lang.IllegalAccessException"
}
