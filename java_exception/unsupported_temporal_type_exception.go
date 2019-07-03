package java_exception

type UnsupportedTemporalTypeException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewUnsupportedTemporalTypeException(detailMessage string) *UnsupportedTemporalTypeException {
	return &UnsupportedTemporalTypeException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e UnsupportedTemporalTypeException) Error() string {
	return e.DetailMessage
}

func (UnsupportedTemporalTypeException) JavaClassName() string {
	return "java.time.temporal.UnsupportedTemporalTypeException"
}
