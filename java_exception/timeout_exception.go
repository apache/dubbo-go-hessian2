package java_exception

type TimeoutException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewTimeoutException(detailMessage string) *TimeoutException {
	return &TimeoutException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e TimeoutException) Error() string {
	return e.DetailMessage
}

func (TimeoutException) JavaClassName() string {
	return "java.util.concurrent.TimeoutException"
}
