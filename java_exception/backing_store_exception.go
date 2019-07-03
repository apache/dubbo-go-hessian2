package java_exception

type BackingStoreException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewBackingStoreException(detailMessage string) *BackingStoreException {
	return &BackingStoreException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e BackingStoreException) Error() string {
	return e.DetailMessage
}

func (BackingStoreException) JavaClassName() string {
	return "java.util.prefs.BackingStoreException"
}
