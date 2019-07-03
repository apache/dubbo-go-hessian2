package java_exception

type ExecutionException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewExecutionException(detailMessage string) *ExecutionException {
	return &ExecutionException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e ExecutionException) Error() string {
	return e.DetailMessage
}

func (ExecutionException) JavaClassName() string {
	return "java.util.concurrent.ExecutionException"
}
