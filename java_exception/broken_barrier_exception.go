package java_exception

type BrokenBarrierException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewBrokenBarrierException(detailMessage string) *BrokenBarrierException {
	return &BrokenBarrierException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e BrokenBarrierException) Error() string {
	return e.DetailMessage
}

func (BrokenBarrierException) JavaClassName() string {
	return "java.util.concurrent.BrokenBarrierException"
}
