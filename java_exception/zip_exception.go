package java_exception

type ZipException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewZipException(detailMessage string) *ZipException {
	return &ZipException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e ZipException) Error() string {
	return e.DetailMessage
}

func (ZipException) JavaClassName() string {
	return "java.util.zip.ZipException"
}
