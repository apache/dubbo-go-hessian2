package java_sql_time

import "time"

type Date struct {
	time.Time
}

func (d *Date) GetTime() time.Time {
	return d.Time
}

func (d *Date) SetTime(time time.Time) {
	d.Time = time
}

func (Date) JavaClassName() string {
	return "java.sql.Date"
}

func (d *Date) ValueOf(dateStr string) error {
	time, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	d.Time = time
	return nil
}

// nolint
func (d *Date) Year() int {
	return d.Time.Year()
}

// nolint
func (d *Date) Month() time.Month {
	return d.Time.Month()
}

// nolint
func (d *Date) Day() int {
	return d.Time.Day()
}
