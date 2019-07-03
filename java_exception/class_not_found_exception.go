package java_exception

type ClassNotFoundException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	Ex                   Throwabler
}

func NewClassNotFoundException(detailMessage string, ex Throwabler) *ClassNotFoundException {
	return &ClassNotFoundException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}, Ex: ex}
}

func (e ClassNotFoundException) Error() string {
	return e.DetailMessage
}

func (ClassNotFoundException) JavaClassName() string {
	return "java.lang.ClassNotFoundException"
}
