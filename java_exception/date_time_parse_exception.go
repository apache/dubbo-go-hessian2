package java_exception

type DateTimeParseException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	ParsedString         string
	ErrorIndex           int32
}

func NewDateTimeParseException(detailMessage string, parsedString string, errorIndex int32) *DateTimeParseException {
	return &DateTimeParseException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}, ParsedString: parsedString, ErrorIndex: errorIndex}
}

func (e DateTimeParseException) Error() string {
	return e.DetailMessage
}

func (DateTimeParseException) JavaClassName() string {
	return "java.time.format.DateTimeParseException"
}
