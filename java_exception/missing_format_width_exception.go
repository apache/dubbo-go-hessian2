package java_exception

type MissingFormatWidthException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	S                    string
}

func NewMissingFormatWidthException(s string) *MissingFormatWidthException {
	return &MissingFormatWidthException{S: s, StackTrace: []StackTraceElement{}}
}

func (e MissingFormatWidthException) Error() string {
	return e.S
}

func (MissingFormatWidthException) JavaClassName() string {
	return "java.util.MissingFormatWidthException"
}
