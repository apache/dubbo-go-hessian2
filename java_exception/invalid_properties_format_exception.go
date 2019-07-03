package java_exception

type InvalidPropertiesFormatException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewInvalidPropertiesFormatException(detailMessage string) *InvalidPropertiesFormatException {
	return &InvalidPropertiesFormatException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e InvalidPropertiesFormatException) Error() string {
	return e.DetailMessage
}

func (InvalidPropertiesFormatException) JavaClassName() string {
	return "java.util.InvalidPropertiesFormatException"
}
