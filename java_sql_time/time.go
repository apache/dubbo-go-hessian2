package java_sql_time

import "time"

type Time struct {
	time.Time
}

func (Time) JavaClassName() string {
	return "java.sql.Time"
}

func (t *Time) GetTime() time.Time {
	return t.Time
}

// nolint
func (t *Time) Hour() int {
	return t.Time.Hour()
}

// nolint
func (t *Time) Minute() int {
	return t.Time.Minute()
}

// nolint
func (t *Time) Second() int {
	return t.Time.Second()
}

func (t *Time) SetTime(time time.Time) {
	t.Time = time
}

func (t *Time) ValueOf(timeStr string) error {
	time, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return err
	}
	t.Time = time
	return nil
}
