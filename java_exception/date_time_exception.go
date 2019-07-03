package java_exception

type DateTimeException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewDateTimeException(detailMessage string) *DateTimeException {
	return &DateTimeException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e DateTimeException) Error() string {
	return e.DetailMessage
}

func (DateTimeException) JavaClassName() string {
	return "java.time.DateTimeException"
}
