package java_exception

type JarException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewJarException(detailMessage string) *JarException {
	return &JarException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e JarException) Error() string {
	return e.DetailMessage
}

func (JarException) JavaClassName() string {
	return "java.util.jar.JarException"
}
