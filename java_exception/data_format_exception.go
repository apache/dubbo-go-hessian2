package java_exception

type DataFormatException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewDataFormatException(detailMessage string) *DataFormatException {
	return &DataFormatException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e DataFormatException) Error() string {
	return e.DetailMessage
}

func (DataFormatException) JavaClassName() string {
	return "java.util.zip.DataFormatException"
}
