package util

import (
	"time"
)

//Dawn returns today dawn time(00:00) in current time zone.
func Dawn() *time.Time {
	now := time.Now()
	t := now.Round(24 * time.Hour)
	if t.After(now) {
		t = t.AddDate(0, 0, -1)
	}
	t = t.Add(-time.Hour * time.Duration(t.Hour()))
	return &t
}
