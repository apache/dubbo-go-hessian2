package java_sql_time

import "time"

type JavaSqlTime interface {
	// ValueOf parse time string which format likes '2006-01-02 15:04:05'
	ValueOf(timeStr string) error
	// SetTime for decode time
	SetTime(time time.Time)
	JavaClassName() string
	// GetTime used to time
	GetTime() time.Time
}
