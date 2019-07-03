package java_exception

type IllegalClassFormatException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewIllegalClassFormatException(detailMessage string) *IllegalClassFormatException {
	return &IllegalClassFormatException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e IllegalClassFormatException) Error() string {
	return e.DetailMessage
}

func (IllegalClassFormatException) JavaClassName() string {
	return "java.lang.instrument.IllegalClassFormatException"
}
