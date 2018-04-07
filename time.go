package gutil

import (
	"time"
)

//Dawn returns today dawn time(00:00) in current time zone.
func Dawn() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
