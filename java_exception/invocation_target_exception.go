package java_exception

type InvocationTargetException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	Target               Throwabler
}

func NewInvocationTargetException(target Throwabler, detailMessage string) *InvocationTargetException {
	return &InvocationTargetException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}, Target: target}
}

func (e InvocationTargetException) Error() string {
	return e.DetailMessage
}

func (InvocationTargetException) JavaClassName() string {
	return "java.lang.reflect.InvocationTargetException"
}
