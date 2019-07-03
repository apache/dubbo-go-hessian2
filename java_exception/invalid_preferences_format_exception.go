package java_exception

type InvalidPreferencesFormatException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewInvalidPreferencesFormatException(detailMessage string) *InvalidPreferencesFormatException {
	return &InvalidPreferencesFormatException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e InvalidPreferencesFormatException) Error() string {
	return e.DetailMessage
}

func (InvalidPreferencesFormatException) JavaClassName() string {
	return "java.util.prefs.InvalidPreferencesFormatException"
}
