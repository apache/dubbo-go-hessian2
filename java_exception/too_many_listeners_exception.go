package java_exception

type TooManyListenersException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewTooManyListenersException(detailMessage string) *TooManyListenersException {
	return &TooManyListenersException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e TooManyListenersException) Error() string {
	return e.DetailMessage
}

func (TooManyListenersException) JavaClassName() string {
	return "java.util.TooManyListenersException"
}
