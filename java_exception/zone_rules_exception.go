package java_exception

type ZoneRulesException struct {
	SerialVersionUID     int64
	DetailMessage        string
	SuppressedExceptions []Throwabler
	StackTrace           []StackTraceElement
	Cause                Throwabler
}

func NewZoneRulesException(detailMessage string) *ZoneRulesException {
	return &ZoneRulesException{DetailMessage: detailMessage, StackTrace: []StackTraceElement{}}
}

func (e ZoneRulesException) Error() string {
	return e.DetailMessage
}

func (ZoneRulesException) JavaClassName() string {
	return "java.time.zone.ZoneRulesException"
}
