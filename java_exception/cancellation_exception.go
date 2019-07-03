package java_exception

type CancellationException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewCancellationException(detailMessage string) *CancellationException {
	return &CancellationException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e CancellationException) Error() string {
	return e.DetailMessage
}

func (CancellationException) JavaClassName() string {
	return "java.util.concurrent.CancellationException"
}
