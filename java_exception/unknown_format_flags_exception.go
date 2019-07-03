package java_exception

type UnknownFormatFlagsException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	Flags                string
}

func NewUnknownFormatFlagsException(flags string) *UnknownFormatFlagsException {
	return &UnknownFormatFlagsException{Flags: flags, StackTrace: []StackTraceElement{}}
}

func (e UnknownFormatFlagsException) Error() string {
	return "Flags = " + e.Flags
}

func (UnknownFormatFlagsException) JavaClassName() string {
	return "java.util.UnknownFormatFlagsException"
}
