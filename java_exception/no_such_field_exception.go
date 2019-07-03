package java_exception

type NoSuchFieldException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewNoSuchFieldException(detailMessage string) *NoSuchFieldException {
	return &NoSuchFieldException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e NoSuchFieldException) Error() string {
	return e.DetailMessage
}

func (NoSuchFieldException) JavaClassName() string {
	return "java.lang.NoSuchFieldException"
}
