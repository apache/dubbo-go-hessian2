package java_exception

import "fmt"

type IllegalFormatCodePointException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	C                    int32
}

func NewIllegalFormatCodePointException(c int32) *IllegalFormatCodePointException {
	return &IllegalFormatCodePointException{C: c, StackTrace: []StackTraceElement{}}
}

func (e IllegalFormatCodePointException) Error() string {
	return fmt.Sprintf("Code point = %#x", e.C)
}

func (IllegalFormatCodePointException) JavaClassName() string {
	return "java.util.IllegalFormatCodePointException"
}
