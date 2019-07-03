package java_exception

type NoSuchMethodException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewNoSuchMethodException(detailMessage string) *NoSuchMethodException {
	return &NoSuchMethodException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e NoSuchMethodException) Error() string {
	return e.DetailMessage
}

func (NoSuchMethodException) JavaClassName() string {
	return "java.lang.NoSuchMethodException"
}
