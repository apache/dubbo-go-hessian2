package java_exception

type InstantiationException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewInstantiationException(detailMessage string) *InstantiationException {
	return &InstantiationException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e InstantiationException) Error() string {
	return e.DetailMessage
}

func (InstantiationException) JavaClassName() string {
	return "java.lang.InstantiationException"
}
