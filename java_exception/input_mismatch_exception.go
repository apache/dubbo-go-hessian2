package java_exception

type InputMismatchException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewInputMismatchException(detailMessage string) *InputMismatchException {
	return &InputMismatchException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e InputMismatchException) Error() string {
	return e.DetailMessage
}

func (InputMismatchException) JavaClassName() string {
	return "java.util.InputMismatchException"
}
