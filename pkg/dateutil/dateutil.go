package dateutil

import "time"

const ISO8601 = "2006-01-02"

func TruncateToMidnight(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func SetYearMonth(t *time.Time, y int, m time.Month) time.Time {
	return time.Date(y, m, t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
