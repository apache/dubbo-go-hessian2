package java_exception

import "strconv"

type IllegalFormatPrecisionException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	P                    int32
}

func NewIllegalFormatPrecisionException(p int32) *IllegalFormatPrecisionException {
	return &IllegalFormatPrecisionException{P: p, StackTrace: []StackTraceElement{}}
}

func (e IllegalFormatPrecisionException) Error() string {
	return strconv.Itoa(int(e.P))
}

func (IllegalFormatPrecisionException) JavaClassName() string {
	return "java.util.IllegalFormatPrecisionException"
}
