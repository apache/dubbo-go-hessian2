package java_exception

import "fmt"

type MissingFormatArgumentException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	S                    string
}

func NewMissingFormatArgumentException(s string) *MissingFormatArgumentException {
	return &MissingFormatArgumentException{S: s, StackTrace: []StackTraceElement{}}
}

func (e MissingFormatArgumentException) Error() string {
	return fmt.Sprintf("Format specifier '%s'", e.S)
}

func (MissingFormatArgumentException) JavaClassName() string {
	return "java.util.MissingFormatArgumentException"
}
