package java_exception

import "fmt"

type UnknownFormatConversionException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
	S                    string
}

func NewUnknownFormatConversionException(s string) *UnknownFormatConversionException {
	return &UnknownFormatConversionException{S: s, StackTrace: []StackTraceElement{}}
}

func (e UnknownFormatConversionException) Error() string {
	return fmt.Sprintf("Conversion = '%s'", e.S)
}

func (UnknownFormatConversionException) JavaClassName() string {
	return "java.util.UnknownFormatConversionException"
}
