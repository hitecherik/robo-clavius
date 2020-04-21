package dateutil

import "time"

const (
	Day            = time.Hour * 24
	ISO8601 string = "2006-01-02"
)

func TruncateToMidnight(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
