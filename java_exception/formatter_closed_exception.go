package java_exception

type FormatterClosedException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewFormatterClosedException() *FormatterClosedException {
	return &FormatterClosedException{StackTrace: []StackTraceElement{}}
}

func (e FormatterClosedException) Error() string {
	return e.DetailMessage
}

func (FormatterClosedException) JavaClassName() string {
	return "java.util.FormatterClosedException"
}
