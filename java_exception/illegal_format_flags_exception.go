package java_exception

import "fmt"

type IllegalFormatFlagsException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	Flags                string
}

func NewIllegalFormatFlagsException(flags string) *IllegalFormatFlagsException {
	return &IllegalFormatFlagsException{Flags: flags, StackTrace: []StackTraceElement{}}
}

func (e IllegalFormatFlagsException) Error() string {
	return fmt.Sprintf("Flags = '%s'", e.Flags)
}

func (IllegalFormatFlagsException) JavaClassName() string {
	return "java.util.IllegalFormatFlagsException"
}
